package repository

import (
	"context"

	"github.com/Rolan335/project/internal/model"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../mocks/blogrepository.go -package=mocks
type BlogRepository interface {
	GetBlog(ctx context.Context, blogID uuid.UUID) (model.DbBlog, error)
	AddBlog(ctx context.Context, blog model.DbBlog) (uuid.UUID, error)
	UpdateBlog(ctx context.Context, blog model.DbBlog) (model.DbBlog, error)
	DeleteBlog(ctx context.Context, blogID uuid.UUID) error
	GetPost(ctx context.Context, postID uuid.UUID) (model.DbPost, error)
	GetPosts(ctx context.Context, BlogID uuid.UUID) ([]model.DbPost, error)
	AddPost(ctx context.Context, post model.DbPost) (uuid.UUID, error)
	UpdatePost(ctx context.Context, post model.DbPost) (model.DbPost, error)
	DeletePost(ctx context.Context, postID uuid.UUID, blogID uuid.UUID) error
}
