package usecase

import (
	"context"

	"github.com/Rolan335/project/internal/model"
)

type BlogUsecaseInterface interface {
	GetBlog(ctx context.Context, req model.BlogGetReq) (model.BlogGetResp, error)
	AddBlog(ctx context.Context, req model.BlogPostReq) (model.BlogPostResp, error)
	UpdateBlog(ctx context.Context, req model.BlogPutReq) (model.BlogPutResp, error)
	DeleteBlog(ctx context.Context, req model.BlogDeleteReq) error
	GetPost(ctx context.Context, req model.PostGetReq) (model.PostGetResp, error)
	GetPosts(ctx context.Context, req model.PostsGetReq) ([]model.PostGetResp, error)
	AddPost(ctx context.Context, req model.PostPostReq) (model.PostPostResp, error)
	UpdatePost(ctx context.Context, req model.PostPutReq) (model.PostPutResp, error)
	DeletePost(ctx context.Context, req model.PostDeleteReq) error
}
