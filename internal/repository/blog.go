package repository

import (
	"context"

	"github.com/Rolan335/project/internal/apperror"
	"github.com/Rolan335/project/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

type BlogRepo struct {
	db *pgxpool.Pool
}

func NewBlogRepo(conn *pgxpool.Pool) *BlogRepo {
	return &BlogRepo{
		db: conn,
	}
}

func (r *BlogRepo) GetBlog(ctx context.Context, blogID uuid.UUID) (model.DbBlog, error) {
	tracer := otel.Tracer("project")
	ctx, span := tracer.Start(ctx, "DB")
	defer span.End()

	var blog model.DbBlog
	if err := pgxscan.Get(ctx, r.db, &blog, "SELECT id, users_id, name, created_at FROM blogs WHERE id = $1", blogID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DbBlog{}, apperror.ErrNotFound
		}
		return model.DbBlog{}, errors.Wrap(err, "blogprovider.BlogRepo.GetBlog")
	}
	return blog, nil
}

func (r *BlogRepo) AddBlog(ctx context.Context, blog model.DbBlog) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "blogprovider.BlogRepo.AddBlog")
	}
	//В идеале потом перекинуть это в отдельный метод для реги юзера
	UserID := blog.UserID
	if _, err := tx.Exec(ctx, "INSERT INTO users(id) values($1) ON CONFLICT (id) DO NOTHING", UserID); err != nil {
		tx.Rollback(ctx)
		return uuid.Nil, errors.Wrap(err, "blogprovider.BlogRepo.AddBlog")
	}

	blogID := blog.ID
	_, err = tx.Exec(ctx, "INSERT INTO blogs(id, users_id, name, created_at) values($1, $2, $3, $4)",
		blog.ID,
		blog.UserID,
		blog.Name,
		blog.CreatedAt,
	)
	if err != nil {
		tx.Rollback(ctx)
		return uuid.Nil, errors.Wrap(err, "blogprovider.BlogRepo.AddBlog")
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return uuid.Nil, errors.Wrap(err, "blogprovider.BlogRepo.AddBlog")
	}
	return blogID, nil
}

func (r *BlogRepo) UpdateBlog(ctx context.Context, blog model.DbBlog) (model.DbBlog, error) {
	var blogRes model.DbBlog
	query := "UPDATE blogs SET users_id = $1, name = $2 WHERE id = $3 RETURNING id, users_id, name, created_at"
	if err := pgxscan.Get(ctx, r.db, &blogRes, query, blog.UserID, blog.Name, blog.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DbBlog{}, apperror.ErrNotFound
		}
		return model.DbBlog{}, errors.Wrap(err, "blogprovider.BlogRepo.UpdateBlog")
	}
	return blogRes, nil
}

func (r *BlogRepo) DeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "blogprovider.BlogRepo.DeleteBlog")
	}
	defer tx.Rollback(ctx)
	cmdTag, err := tx.Exec(ctx, "DELETE FROM blogs WHERE id = $1", blogID)
	if err != nil {
		return errors.Wrap(err, "blogprovider.BlogRepo.DeleteBlog")
	}
	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	//deleting all posts in deleted blog
	if _, err := tx.Exec(ctx, "DELETE FROM posts WHERE blog_id = $1", blogID); err != nil {
		return errors.Wrap(err, "blogprovider.BlogRepo.DeleteBlog")
	}
	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "blogprovider.BlogRepo.DeleteBlog")
	}
	return nil
}

func (r *BlogRepo) GetPost(ctx context.Context, postID uuid.UUID) (model.DbPost, error) {
	var post model.DbPost
	if err := pgxscan.Get(ctx, r.db, &post, "SELECT id, blogs_id, title, text, created_at FROM posts WHERE id = $1", postID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DbPost{}, apperror.ErrNotFound
		}
		return model.DbPost{}, errors.Wrap(err, "blogprovider.BlogRepo.GetPost")
	}
	return post, nil
}

func (r *BlogRepo) GetPosts(ctx context.Context, BlogID uuid.UUID) ([]model.DbPost, error) {
	var posts []model.DbPost
	if err := pgxscan.Select(ctx, r.db, &posts, "SELECT id, blogs_id, title, text, created_at FROM posts WHERE blogs_id = $1", BlogID); err != nil {
		return nil, errors.Wrap(err, "blogprovider.BlogRepo.GetPosts")
	}
	return posts, nil
}

func (r *BlogRepo) AddPost(ctx context.Context, post model.DbPost) (uuid.UUID, error) {
	if err := r.db.QueryRow(ctx, "SELECT id FROM blogs WHERE id = $1", post.BlogID).Scan(nil); err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, apperror.ErrNotFound
		}
		return uuid.Nil, errors.Wrap(err, "blogprovider.BlogRepo.AddPost")
	}
	_, err := r.db.Exec(ctx, "INSERT INTO posts(id, blogs_id, title, text, created_at) VALUES($1, $2, $3, $4, $5)",
		post.ID,
		post.BlogID,
		post.Title,
		post.Text,
		post.CreatedAt,
	)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "blogprovider.BlogRepo.AddPost")
	}

	return post.ID, nil
}

func (r *BlogRepo) UpdatePost(ctx context.Context, post model.DbPost) (model.DbPost, error) {
	var postRes model.DbPost
	// Обновляем только title и text
	query := "UPDATE posts SET title = $1, text = $2 WHERE id = $3 AND blogs_id = $4 RETURNING id, blogs_id, title, text, created_at"
	err := pgxscan.Get(ctx, r.db, &postRes, query, post.Title, post.Text, post.ID, post.BlogID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DbPost{}, apperror.ErrNotFound
		}
		return model.DbPost{}, errors.Wrap(err, "blogprovider.BlogRepo.UpdatePost")
	}
	return postRes, nil
}

func (r *BlogRepo) DeletePost(ctx context.Context, postID uuid.UUID, blogID uuid.UUID) error {
	cmdTag, err := r.db.Exec(ctx, "DELETE FROM posts WHERE id = $1 AND blogs_id = $2", postID, blogID)
	if err != nil {
		return errors.Wrap(err, "blogprovider.BlogRepo.DeletePost")
	}
	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
