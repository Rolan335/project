package handler

import (
	"errors"

	"github.com/Rolan335/project/internal/apperr"
	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

const (
	BlogIDParam = "blog_id"
	PostIDParam = "post_id"
)

type Handler struct {
	validate *validator.Validate
	usecase  usecase.BlogUsecase
}

func New(usecase usecase.BlogUsecase, validate *validator.Validate) *Handler {
	return &Handler{
		validate: validate,
		usecase:  usecase,
	}
}

func (h *Handler) GetBlog(c *fiber.Ctx) error {
	ctx := c.UserContext()
	tracer := otel.Tracer("project")
	_, span := tracer.Start(ctx, "GetBlog")
	defer span.End()

	blogID, err := uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	blog, err := h.usecase.GetBlog(ctx, model.BlogGetReq{BlogID: blogID})
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.JSON(blog)
}

func (h *Handler) CreateBlog(c *fiber.Ctx) error {
	var blog model.BlogPostReq
	if err := c.BodyParser(&blog); err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.validate.Struct(blog); err != nil {
		log.Err(err).Msg("")
		return fiber.ErrBadRequest
	}
	id, err := h.usecase.AddBlog(c.Context(), blog)
	if err != nil {
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.JSON(id)
}

func (h *Handler) UpdateBlog(c *fiber.Ctx) error {
	var blog model.BlogPutReq
	var err error
	blog.BlogID, err = uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := c.BodyParser(&blog); err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.validate.Struct(blog); err != nil {
		log.Err(err).Msg("")
		return fiber.ErrBadRequest
	}
	resp, err := h.usecase.UpdateBlog(c.Context(), blog)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.JSON(resp)
}

func (h *Handler) DeleteBlog(c *fiber.Ctx) error {
	var blog model.BlogDeleteReq
	var err error
	blog.BlogID, err = uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.usecase.DeleteBlog(c.Context(), blog); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetPosts(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	posts, err := h.usecase.GetPosts(c.Context(), model.PostsGetReq{BlogID: blogID})
	if err != nil {
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	if len(posts) == 0 {
		return fiber.ErrNotFound
	}
	return c.JSON(posts)
}

// TODO: Не проверяет валидность blog_id. Т.е. запрос api/blog/левый-uuid/posts/валидный-uuid вернёт пост.
func (h *Handler) GetPost(c *fiber.Ctx) error {
	var req model.PostGetReq
	var err error
	req.BlogID, err = uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.PostID, err = uuid.Parse(c.Params(PostIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	post, err := h.usecase.GetPost(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}

	return c.JSON(post)
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	var req model.PostPostReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	var err error
	req.BlogID, err = uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.validate.Struct(req); err != nil {
		log.Err(err).Msg("")
		return fiber.ErrBadRequest
	}
	resp, err := h.usecase.AddPost(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.JSON(resp)
}

func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	var req model.PostPutReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	var err error
	req.BlogID, err = uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.PostID, err = uuid.Parse(c.Params(PostIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.validate.Struct(req); err != nil {
		log.Err(err).Msg("")
		return fiber.ErrBadRequest
	}
	resp, err := h.usecase.UpdatePost(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.JSON(resp)
}

func (h *Handler) DeletePost(c *fiber.Ctx) error {
	var req model.PostDeleteReq
	var err error
	req.BlogID, err = uuid.Parse(c.Params(BlogIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.PostID, err = uuid.Parse(c.Params(PostIDParam))
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err = h.usecase.DeletePost(c.Context(), req); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return fiber.ErrNotFound
		}
		log.Err(err).Msg("")
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusOK)
}
