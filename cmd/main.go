package main

import (
	"github.com/Rolan335/project/internal/handler"
	"github.com/Rolan335/project/internal/repo/inmemory"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	api := app.Group("/api")

	repository := inmemory.New()

	h := handler.New(repository)

	api.Get("/blog/:blog_id", h.GetBlog) // blog_id: uuid
	api.Post("/blog", h.CreateBlog) // reqBody {"user_id": "abc123", "name": "firstBlog"} | res {"blog_id":"aaaaa-aasfsf-43ife-fdfs-fddd"}
	api.Put("/blog/:blog_id", h.UpdateBlog) // blog_id: uuid | reqBody {"user_id": "abc123", "name": "firstBlog"} res model.Blog
	api.Delete("/blog/:blog_id", h.DeleteBlog)
	api.Get("/blog/:blog_id/posts", h.GetPosts)
	api.Get("/blog/:blog_id/posts/:post_id", h.GetPost)
	api.Put("/blog/:blog_id/posts/:post_id", h.UpdatePost)
	api.Post("/blog/:blog_id/posts", h.CreatePost) //reqBody {"title": "first Post", "text": "first post textttt"} | res {"post_id":"aaaaa-aasfsf-43ife-fdfs-fddd"}
	api.Delete("/blog/:blog_id/posts/:post_id", h.DeletePost)

	app.Listen(":8080")
}
