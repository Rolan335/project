package repository

import (
	"context"

	"github.com/Rolan335/project/internal/model/dto"
	"github.com/google/uuid"
)

type BlogRepoInterface interface {
	GetBlog(ctx context.Context, blogID uuid.UUID) (dto.DbBlog, error)
	AddBlog(ctx context.Context, blog dto.DbBlog) (uuid.UUID, error)
	UpdateBlog(ctx context.Context, blog dto.DbBlog) (dto.DbBlog, error)
	DeleteBlog(ctx context.Context, blogID uuid.UUID) error
	GetPost(ctx context.Context, postID uuid.UUID) (dto.DbPost, error)
	GetPosts(ctx context.Context, BlogID uuid.UUID) ([]dto.DbPost, error)
	AddPost(ctx context.Context, post dto.DbPost) (uuid.UUID, error)
	UpdatePost(ctx context.Context, post dto.DbPost) (dto.DbPost, error)
	DeletePost(ctx context.Context, postID uuid.UUID, blogID uuid.UUID) error
}
