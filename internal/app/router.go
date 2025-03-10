package app

import (
	"github.com/Rolan335/project/internal/handler"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

func GetRouter(handle *handler.Handler) *fiber.App {
	app := fiber.New()
	api := app.Group("/api")
	api.Use(otelfiber.Middleware())

	api.Get("/blog/:blog_id", handle.GetBlog)
	api.Post("/blog", handle.CreateBlog)
	api.Put("/blog/:blog_id", handle.UpdateBlog)
	api.Delete("/blog/:blog_id", handle.DeleteBlog)
	api.Get("/blog/:blog_id/posts", handle.GetPosts)
	api.Get("/blog/:blog_id/posts/:post_id", handle.GetPost)
	api.Put("/blog/:blog_id/posts/:post_id", handle.UpdatePost)
	api.Post("/blog/:blog_id/posts", handle.CreatePost)
	api.Delete("/blog/:blog_id/posts/:post_id", handle.DeletePost)

	return app
}
