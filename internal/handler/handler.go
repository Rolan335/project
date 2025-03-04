package handler

import (
	"context"
	"errors"
	"time"

	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	BlogIDParam = "blog_id"
	PostIDParam = "post_id"
)

type Repo interface {
	GetBlog(ctx context.Context, blogID string) (model.Blog, error)
	AddBlog(ctx context.Context, blog model.Blog) (string, error)
	UpdateBlog(ctx context.Context, blog model.Blog) (model.Blog, error)
	DeleteBlog(ctx context.Context, blogID string) error
	GetPost(ctx context.Context, postID string) (model.Post, error)
	GetPosts(ctx context.Context, BlogID string) ([]model.Post, error)
	AddPost(ctx context.Context, post model.Post) (string, error)
	UpdatePost(ctx context.Context, post model.Post) (model.Post, error)
	DeletePost(ctx context.Context, postID string) error
}

type Handler struct {
	r Repo
}

func New(repo Repo) *Handler {
	return &Handler{
		r: repo,
	}
}

func (h *Handler) GetBlog(c *fiber.Ctx) error {
	blog, err := h.r.GetBlog(c.Context(), c.AllParams()[BlogIDParam])
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
	}
	return c.JSON(blog)
}

func (h *Handler) CreateBlog(c *fiber.Ctx) error {
	//TODO: норм валидация
	var blog model.Blog
	if err := c.BodyParser(&blog); err != nil {
		return fiber.ErrBadRequest
	}
	//Перенести в service
	uuid, _ := uuid.NewRandom()
	blog.ID = uuid.String()
	blog.CreatedAt = time.Now()
	id, err := h.r.AddBlog(c.Context(), blog)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(fiber.Map{"blog_id": id})
}

func (h *Handler) UpdateBlog(c *fiber.Ctx) error {
	var blog model.Blog
	if err := c.BodyParser(&blog); err != nil {
		return fiber.ErrBadRequest
	}
	blogUpd, err := h.r.UpdateBlog(c.Context(), blog)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}
	return c.JSON(blogUpd)
}

func (h *Handler) DeleteBlog(c *fiber.Ctx) error {
	err := h.r.DeleteBlog(c.Context(), c.AllParams()[BlogIDParam])
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(200)
}

func (h *Handler) GetPosts(c *fiber.Ctx) error {
	posts, err := h.r.GetPosts(c.Context(), c.AllParams()[BlogIDParam])
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(posts)
}

func (h *Handler) GetPost(c *fiber.Ctx) error {
	post, err := h.r.GetPost(c.Context(), c.AllParams()[PostIDParam])
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}
	return c.JSON(post)
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	//TODO: норм валидация
	var post model.Post
	if err := c.BodyParser(&post); err != nil {
		return fiber.ErrBadRequest
	}
	post.BlogID = c.AllParams()[BlogIDParam]
	uuid, _ := uuid.NewRandom()
	post.ID = uuid.String()
	post.CreatedAt = time.Now()
	id, err := h.r.AddPost(c.Context(), post)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(fiber.Map{"post_id": id})
}

func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	var post model.Post
	if err := c.BodyParser(&post); err != nil {
		return fiber.ErrBadRequest
	}
	postUpd, err := h.r.UpdatePost(c.Context(), post)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}
	return c.JSON(postUpd)
}

func (h *Handler) DeletePost(c *fiber.Ctx) error {
	err := h.r.DeletePost(c.Context(), c.AllParams()[PostIDParam])
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(200)
}
