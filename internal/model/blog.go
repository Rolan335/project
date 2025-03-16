package model

import (
	"time"

	"github.com/google/uuid"
)

type BlogGetReq struct {
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
}

type BlogGetResp struct {
	BlogID    uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
type BlogPostReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
	Name   string    `json:"name" validate:"required,min=1,max=64"`
}

type BlogPostResp struct {
	BlogID uuid.UUID `json:"id"`
}

type BlogPutReq struct {
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
	Name   string    `json:"name" validate:"required,min=1,max=64"`
}

type BlogPutResp struct {
	BlogID    uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogDeleteReq struct {
	BlogID uuid.UUID `json:"id" validate:"required,uuid"`
}

type PostsGetReq struct {
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
}

type PostGetReq struct {
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
	PostID uuid.UUID `json:"post_id" validate:"required,uuid"`
}

type PostGetResp struct {
	PostID    uuid.UUID `json:"post_id"`
	BlogID    uuid.UUID `json:"blog_id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type PostPostReq struct {
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
	Title  string    `json:"title" validate:"required,min=1,max=64"`
	Text   string    `json:"text" validate:"required,min=1,max=2048"`
}

type PostPostResp struct {
	PostID uuid.UUID `json:"post_id"`
}

type PostPutReq struct {
	PostID uuid.UUID `json:"post_id" validate:"required,uuid"`
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
	Title  string    `json:"title" validate:"required,min=1,max=64"`
	Text   string    `json:"text" validate:"required,min=1,max=2048"`
}

type PostPutResp struct {
	PostID    uuid.UUID `json:"post_id"`
	BlogID    uuid.UUID `json:"blog_id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type PostDeleteReq struct {
	PostID uuid.UUID `json:"post_id" validate:"required,uuid"`
	BlogID uuid.UUID `json:"blog_id" validate:"required,uuid"`
}
