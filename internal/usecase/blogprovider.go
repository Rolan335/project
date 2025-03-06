package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/model/dto"
	"github.com/Rolan335/project/internal/repository"
	"github.com/google/uuid"
)

type BlogProvider struct {
	repository repository.BlogRepoInterface
}

func NewBlogProvider(repository repository.BlogRepoInterface) *BlogProvider {
	return &BlogProvider{
		repository: repository,
	}
}

func (b *BlogProvider) GetBlog(ctx context.Context, req model.BlogGetReq) (model.BlogGetResp, error) {
	blogDB, err := b.repository.GetBlog(ctx, req.BlogID)
	if err != nil {
		return model.BlogGetResp{}, fmt.Errorf("usercase.BlogProvider.GetBlog: %w", err)
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
	blogDB := dto.DbBlog{
		ID:        id,
		UserID:    req.UserID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	blogid, err := b.repository.AddBlog(ctx, blogDB)
	if err != nil {
		return model.BlogPostResp{}, fmt.Errorf("usecase.BlogProvider.AddBlog: %w", err)
	}

	return model.BlogPostResp{BlogID: blogid}, nil
}

func (b *BlogProvider) UpdateBlog(ctx context.Context, req model.BlogPutReq) (model.BlogPutResp, error) {
	//Обновляет userID, Name
	blogDB := dto.DbBlog{
		ID:     req.BlogID,
		UserID: req.UserID,
		Name:   req.Name,
	}
	blog, err := b.repository.UpdateBlog(ctx, blogDB)
	if err != nil {
		return model.BlogPutResp{}, fmt.Errorf("usecase.BlogProvider.UpdateBlog: %w", err)
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
		return fmt.Errorf("usecase.BlogProvider.DeleteBlog: %w", err)
	}
	return nil
}
func (b *BlogProvider) GetPost(ctx context.Context, req model.PostGetReq) (model.PostGetResp, error) {
	post, err := b.repository.GetPost(ctx, req.PostID)
	if err != nil {
		return model.PostGetResp{}, fmt.Errorf("usecase.BlogProvider.GetPost: %w", err)
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
		return nil, fmt.Errorf("usecase.BlogProvider.GetPosts: %w", err)
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
	dbPost := dto.DbPost{
		BlogID: req.BlogID,
		Title:  req.Title,
		Text:   req.Text,
	}
	dbPost.ID, _ = uuid.NewRandom()
	dbPost.CreatedAt = time.Now()
	postID, err := b.repository.AddPost(ctx, dbPost)
	if err != nil {
		return model.PostPostResp{}, fmt.Errorf("usecase.BlogProvider.AddPost: %w", err)
	}
	return model.PostPostResp{PostID: postID}, nil
}
func (b *BlogProvider) UpdatePost(ctx context.Context, req model.PostPutReq) (model.PostPutResp, error) {
	dbPost := dto.DbPost{
		ID:     req.PostID,
		BlogID: req.BlogID,
		Title:  req.Title,
		Text:   req.Text,
	}
	post, err := b.repository.UpdatePost(ctx, dbPost)
	if err != nil {
		return model.PostPutResp{}, fmt.Errorf("usecase.BlogProvider.UpdatePost: %w", err)
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
		return fmt.Errorf("usecase.BlogProvider.DeletePost: %w", err)
	}
	return nil
}
