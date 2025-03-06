package blogprovider

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rolan335/project/internal/model/dto"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{
		db: conn,
	}
}

func (r *Repository) GetBlog(ctx context.Context, blogID uuid.UUID) (dto.DbBlog, error) {
	var blog dto.DbBlog
	if err := pgxscan.Get(ctx, r.db, &blog, "SELECT id, users_id, name, created_at FROM blogs WHERE id = $1", blogID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.DbBlog{}, nil
		}
		return dto.DbBlog{}, fmt.Errorf("blogprovider.Repository.GetBlog: %w", err)
	}
	return blog, nil
}

func (r *Repository) AddBlog(ctx context.Context, blog dto.DbBlog) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("blogprovider.Repository.AddBlog: %w", err)
	}
	//В идеале потом перекинуть это в отдельный метод для реги юзера
	UserID := blog.UserID
	if _, err := tx.Exec(ctx, "INSERT INTO users(id) values($1)", UserID); err != nil {
		tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("blogprovider.Repository.AddBlog: %w", err)
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
		return uuid.Nil, fmt.Errorf("blogprovider.Repository.AddBlog: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("blogprovider.Repository.AddBlog: %w", err)
	}
	return blogID, nil
}

// ??? Тут приходит полная dto для db, но обновить из всех полей мы можем только name и user_id. Как лучше оформить обновление?
func (r *Repository) UpdateBlog(ctx context.Context, blog dto.DbBlog) (dto.DbBlog, error) {
	var blogRes dto.DbBlog
	query := "UPDATE blogs SET users_id = $1, name = $2 WHERE id = $3 RETURNING id, users_id, name, created_at"
	if err := pgxscan.Get(ctx, r.db, &blogRes, query, blog.UserID, blog.Name, blog.ID); err != nil {
		return dto.DbBlog{}, fmt.Errorf("blogprovider.Repository.UpdateBlog: %w", err)
	}
	return blogRes, nil
}

func (r *Repository) DeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	if _, err := r.db.Exec(ctx, "DELETE FROM blogs WHERE id = $1", blogID); err != nil {
		return fmt.Errorf("blogprovider.Repository.DeleteBlog: %w", err)
	}
	return nil
}
func (r *Repository) GetPost(ctx context.Context, postID uuid.UUID) (dto.DbPost, error) {
	var post dto.DbPost
	if err := pgxscan.Get(ctx, r.db, &post, "SELECT id, blogs_id, title, text, created_at FROM posts WHERE id = $1", postID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.DbPost{}, nil
		}
		return dto.DbPost{}, fmt.Errorf("blogprovider.Repository.GetPost: %w", err)
	}
	return post, nil
}
func (r *Repository) GetPosts(ctx context.Context, BlogID uuid.UUID) ([]dto.DbPost, error) {
	var posts []dto.DbPost
	if err := pgxscan.Select(ctx, r.db, &posts, "SELECT id, blogs_id, title, text, created_at FROM posts WHERE blogs_id = $1", BlogID); err != nil {
		return nil, fmt.Errorf("blogprovider.Repository.GetPosts: %w", err)
	}
	return posts, nil
}
func (r *Repository) AddPost(ctx context.Context, post dto.DbPost) (uuid.UUID, error) {
	if err := r.db.QueryRow(ctx, "SELECT id FROM blogs WHERE id = $1", post.BlogID).Scan(nil); err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, nil
		}
		return uuid.Nil, fmt.Errorf("postprovider.Repository.AddPost: %w", err)
	}
	_, err := r.db.Exec(ctx, "INSERT INTO posts(id, blogs_id, title, text, created_at) VALUES($1, $2, $3, $4, $5)",
		post.ID,
		post.BlogID,
		post.Title,
		post.Text,
		post.CreatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("postprovider.Repository.AddPost: %w", err)
	}

	return post.ID, nil
}

// ??? Как лучше сделать update. Норм, если в dto приходит не полная структура? Нужно ли явно задавать какие поля
// ? обновляются или оставить просто query как есть
func (r *Repository) UpdatePost(ctx context.Context, post dto.DbPost) (dto.DbPost, error) {
	var postRes dto.DbPost
	// Обновляем только title и text
	query := "UPDATE posts SET title = $1, text = $2 WHERE id = $3 AND blogs_id = $4 RETURNING id, blogs_id, title, text, created_at"
	err := pgxscan.Get(ctx, r.db, &postRes, query, post.Title, post.Text, post.ID, post.BlogID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.DbPost{}, nil
		}
		return dto.DbPost{}, fmt.Errorf("postprovider.Repository.UpdatePost: %w", err)
	}
	return postRes, nil
}
func (r *Repository) DeletePost(ctx context.Context, postID uuid.UUID, blogID uuid.UUID) error {
	if _, err := r.db.Exec(ctx, "DELETE FROM posts WHERE id = $1 AND blogs_id = $2", postID, blogID); err != nil {
		return fmt.Errorf("postprovider.Repository.DeletePost: %w", err)
	}
	return nil
}
