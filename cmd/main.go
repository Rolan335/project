package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	api := app.Group("/api")
	//Returns blog as {blogID, CreatorID} мб что-то ещё
	api.Get("/blog/:id")
	//Creates a new blog reqBody {creatorID, name} returns {blogID}
	api.Post("/blog")
	//Deletes blog reqBody {creatorID}
	api.Delete("/blog")
	//Returns all posts with blogID as []string{CreatorID, Text, DateTime}
	api.Get("/blog/:id/posts")
	//Creates new post for provided blogID. reqBody {creatorID, Title, Text} returns {PostID}
	api.Post("/blog/:id/posts")
	//Deletes a post reqBody {CreatorID}
	api.Delete("/blog/:blog_id/posts/:post_id")
	//Get comments for provided post returns []model.Comment{}
	api.Get("/blog/:blog_id/posts/:post_id/comments")
	//Create comment for provided post reqBody {CreatorID, Text} return {CommentID}
	api.Post("/blog/:blog_id/posts/:post_id/comments")
	//Deletes a comment reqBody {CreatorID}
	api.Delete("/blog/:blog_id/posts/:post_id/comments/:comment_id")

	//PUT по той же логике.
	
	app.Listen(":8080")
}
