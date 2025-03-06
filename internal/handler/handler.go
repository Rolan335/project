package handler

import (
	"log"

	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	BlogIDParam = "blog_id"
	PostIDParam = "post_id"
)

type Handler struct {
	usecase usecase.BlogUsecaseInterface
}

func New(usecase usecase.BlogUsecaseInterface) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) GetBlog(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	blog, err := h.usecase.GetBlog(c.Context(), model.BlogGetReq{BlogID: blogID})
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if blog == (model.BlogGetResp{}) {
		return fiber.ErrNotFound
	}
	return c.JSON(blog)
}

func (h *Handler) CreateBlog(c *fiber.Ctx) error {
	var blog model.BlogPostReq
	if err := c.BodyParser(&blog); err != nil {
		return fiber.ErrBadRequest
	}
	if err := blog.Validate(); err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}
	id, err := h.usecase.AddBlog(c.Context(), blog)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return c.JSON(id)
}

func (h *Handler) UpdateBlog(c *fiber.Ctx) error {
	var blog model.BlogPutReq
	var err error
	blog.BlogID, err = uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := c.BodyParser(&blog); err != nil {
		return fiber.ErrBadRequest
	}
	if err := blog.Validate(); err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.usecase.UpdateBlog(c.Context(), blog)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return c.JSON(resp)
}

func (h *Handler) DeleteBlog(c *fiber.Ctx) error {
	var blog model.BlogDeleteReq
	var err error
	blog.BlogID, err = uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := blog.Validate(); err != nil {
		return fiber.ErrBadRequest
	}
	//TODO: Тут будет 500 если not found
	if err := h.usecase.DeleteBlog(c.Context(), blog); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(200)
}

// TODO: Не проверяет валидность blog_id. Т.е. запрос api/blog/левый-uuid/posts/валидный-uuid вернёт пост.
func (h *Handler) GetPosts(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	posts, err := h.usecase.GetPosts(c.Context(), model.PostsGetReq{BlogID: blogID})
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if len(posts) == 0 {
		return fiber.ErrNotFound
	}
	return c.JSON(posts)
}

func (h *Handler) GetPost(c *fiber.Ctx) error {
	var req model.PostGetReq
	var err error
	req.BlogID, err = uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.PostID, err = uuid.Parse(c.AllParams()[PostIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	post, err := h.usecase.GetPost(c.Context(), req)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if post == (model.PostGetResp{}) {
		return fiber.ErrNotFound
	}

	return c.JSON(post)
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	var req model.PostPostReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	var err error
	req.BlogID, err = uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.usecase.AddPost(c.Context(), req)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if resp.PostID == uuid.Nil {
		return fiber.ErrNotFound
	}
	return c.JSON(resp)
}

func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	var req model.PostPutReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	var err error
	req.BlogID, err = uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.PostID, err = uuid.Parse(c.AllParams()[PostIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.usecase.UpdatePost(c.Context(), req)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if resp == (model.PostPutResp{}) {
		return fiber.ErrNotFound
	}
	return c.JSON(resp)
}

func (h *Handler) DeletePost(c *fiber.Ctx) error {
	var req model.PostDeleteReq
	var err error
	req.BlogID, err = uuid.Parse(c.AllParams()[BlogIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.PostID, err = uuid.Parse(c.AllParams()[PostIDParam])
	if err != nil {
		return fiber.ErrBadRequest
	}
	//TODO: Тут будет 500 если not found
	err = h.usecase.DeletePost(c.Context(), req)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(200)
}
