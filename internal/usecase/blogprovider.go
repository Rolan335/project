package usecase

import (
	"context"
	"time"

	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

type BlogProvider struct {
	repository repository.BlogRepository
}

func NewBlogProvider(repository repository.BlogRepository) *BlogProvider {
	return &BlogProvider{
		repository: repository,
	}
}

func (b *BlogProvider) GetBlog(ctx context.Context, req model.BlogGetReq) (model.BlogGetResp, error) {
	tracer := otel.Tracer("project")
	_, span := tracer.Start(ctx, "GetBlogUsecase")
	defer span.End()

	blogDB, err := b.repository.GetBlog(ctx, req.BlogID)
	if err != nil {
		return model.BlogGetResp{}, errors.Wrap(err, "usercase.BlogProvider.GetBlog")
	}
	return model.BlogGetResp{
		BlogID:    blogDB.ID,
		UserID:    blogDB.UserID,
		Name:      blogDB.Name,
		CreatedAt: blogDB.CreatedAt,
	}, nil
}
func (b *BlogProvider) AddBlog(ctx context.Context, req model.BlogPostReq) (model.BlogPostResp, error) {
	id, _ := uuid.NewRandom()
	blogDB := model.DbBlog{
		ID:        id,
		UserID:    req.UserID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	blogid, err := b.repository.AddBlog(ctx, blogDB)
	if err != nil {
		return model.BlogPostResp{}, errors.Wrap(err, "usercase.BlogProvider.AddBlog")
	}

	return model.BlogPostResp{BlogID: blogid}, nil
}

func (b *BlogProvider) UpdateBlog(ctx context.Context, req model.BlogPutReq) (model.BlogPutResp, error) {
	//Обновляет userID, Name
	blogDB := model.DbBlog{
		ID:     req.BlogID,
		UserID: req.UserID,
		Name:   req.Name,
	}
	blog, err := b.repository.UpdateBlog(ctx, blogDB)
	if err != nil {
		return model.BlogPutResp{}, errors.Wrap(err, "usercase.BlogProvider.UpdateBlog")
	}

	return model.BlogPutResp{
		BlogID:    blog.ID,
		UserID:    blog.UserID,
		Name:      blog.Name,
		CreatedAt: blog.CreatedAt,
	}, nil
}
func (b *BlogProvider) DeleteBlog(ctx context.Context, req model.BlogDeleteReq) error {
	if err := b.repository.DeleteBlog(ctx, req.BlogID); err != nil {
		return errors.Wrap(err, "usercase.BlogProvider.DeleteBlog")
	}
	return nil
}
func (b *BlogProvider) GetPost(ctx context.Context, req model.PostGetReq) (model.PostGetResp, error) {
	post, err := b.repository.GetPost(ctx, req.PostID)
	if err != nil {
		return model.PostGetResp{}, errors.Wrap(err, "usercase.BlogProvider.GetPost")
	}

	return model.PostGetResp{
		PostID:    post.ID,
		BlogID:    post.BlogID,
		Title:     post.Title,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
	}, nil
}
func (b *BlogProvider) GetPosts(ctx context.Context, req model.PostsGetReq) ([]model.PostGetResp, error) {
	posts, err := b.repository.GetPosts(ctx, req.BlogID)
	if err != nil {
		return nil, errors.Wrap(err, "usercase.BlogProvider.GetPosts")
	}
	resp := make([]model.PostGetResp, 0, len(posts))
	for i := 0; i < len(posts); i++ {
		resp = append(resp, model.PostGetResp{
			PostID:    posts[i].ID,
			BlogID:    posts[i].BlogID,
			Title:     posts[i].Title,
			Text:      posts[i].Text,
			CreatedAt: posts[i].CreatedAt,
		})
	}
	return resp, nil
}
func (b *BlogProvider) AddPost(ctx context.Context, req model.PostPostReq) (model.PostPostResp, error) {
	dbPost := model.DbPost{
		BlogID: req.BlogID,
		Title:  req.Title,
		Text:   req.Text,
	}
	dbPost.ID, _ = uuid.NewRandom()
	dbPost.CreatedAt = time.Now()
	postID, err := b.repository.AddPost(ctx, dbPost)
	if err != nil {
		return model.PostPostResp{}, errors.Wrap(err, "usercase.BlogProvider.AddPost")
	}
	return model.PostPostResp{PostID: postID}, nil
}
func (b *BlogProvider) UpdatePost(ctx context.Context, req model.PostPutReq) (model.PostPutResp, error) {
	dbPost := model.DbPost{
		ID:     req.PostID,
		BlogID: req.BlogID,
		Title:  req.Title,
		Text:   req.Text,
	}
	post, err := b.repository.UpdatePost(ctx, dbPost)
	if err != nil {
		return model.PostPutResp{}, errors.Wrap(err, "usercase.BlogProvider.UpdatePost")
	}
	return model.PostPutResp{
		PostID:    post.ID,
		BlogID:    post.BlogID,
		Title:     post.Title,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
	}, nil
}
func (b *BlogProvider) DeletePost(ctx context.Context, req model.PostDeleteReq) error {
	if err := b.repository.DeletePost(ctx, req.PostID, req.BlogID); err != nil {
		return errors.Wrap(err, "usercase.BlogProvider.DeletePost")
	}
	return nil
}
