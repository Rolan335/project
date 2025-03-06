package app

import (
	"github.com/Rolan335/project/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func GetRouter(handle *handler.Handler) *fiber.App {
	app := fiber.New()
	api := app.Group("/api")

	api.Get("/blog/:blog_id", handle.GetBlog)    // blog_id: uuid
	api.Post("/blog", handle.CreateBlog)         // reqBody {"user_id": "abc123", "name": "firstBlog"} | res {"blog_id":"aaaaa-aasfsf-43ife-fdfs-fddd"}
	api.Put("/blog/:blog_id", handle.UpdateBlog) // blog_id: uuid | reqBody {"user_id": "abc123", "name": "firstBlog"} res model.Blog
	api.Delete("/blog/:blog_id", handle.DeleteBlog)
	api.Get("/blog/:blog_id/posts", handle.GetPosts)
	api.Get("/blog/:blog_id/posts/:post_id", handle.GetPost)
	api.Put("/blog/:blog_id/posts/:post_id", handle.UpdatePost)
	api.Post("/blog/:blog_id/posts", handle.CreatePost) //reqBody {"title": "first Post", "text": "first post textttt"} | res {"post_id":"aaaaa-aasfsf-43ife-fdfs-fddd"}
	api.Delete("/blog/:blog_id/posts/:post_id", handle.DeletePost)

	return app
}
