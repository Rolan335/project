package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type BlogGetReq struct {
	BlogID uuid.UUID `json:"blog_id,omitempty"`
}

type BlogGetResp struct {
	BlogID    uuid.UUID `json:"id,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type BlogPostReq struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
	Name   string    `json:"name,omitempty"`
}

func (b *BlogPostReq) Validate() error {
	if b.UserID == uuid.Nil {
		return errors.New("user_id cannot be empty")
	}
	if b.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

type BlogPostResp struct {
	BlogID uuid.UUID `json:"id,omitempty"`
}

type BlogPutReq struct {
	BlogID uuid.UUID `json:"blog_id,omitempty"`
	UserID uuid.UUID `json:"user_id,omitempty"`
	Name   string    `json:"name,omitempty"`
}

func (b *BlogPutReq) Validate() error {
	if b.BlogID == uuid.Nil {
		return errors.New("blog_id cannot be empty")
	}
	if b.UserID == uuid.Nil {
		return errors.New("user_id cannot be empty")
	}
	if b.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

type BlogPutResp struct {
	BlogID    uuid.UUID `json:"id,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type BlogDeleteReq struct {
	BlogID uuid.UUID `json:"id,omitempty"`
}

func (b *BlogDeleteReq) Validate() error {
	if b.BlogID == uuid.Nil {
		return errors.New("user_id cannot be empty")
	}
	return nil
}

type PostsGetReq struct {
	BlogID uuid.UUID `json:"blog_id,omitempty"`
}

type PostGetReq struct {
	BlogID uuid.UUID `json:"blog_id,omitempty"`
	PostID uuid.UUID `json:"post_id,omitempty"`
}

type PostGetResp struct {
	PostID    uuid.UUID `json:"post_id,omitempty"`
	BlogID    uuid.UUID `json:"blog_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type PostPostReq struct {
	BlogID uuid.UUID `json:"blog_id,omitempty"`
	Title  string    `json:"title,omitempty"`
	Text   string    `json:"text,omitempty"`
}

type PostPostResp struct {
	PostID uuid.UUID `json:"post_id,omitempty"`
}

type PostPutReq struct {
	PostID uuid.UUID `json:"post_id,omitempty"`
	BlogID uuid.UUID `json:"blog_id,omitempty"`
	Title  string    `json:"title,omitempty"`
	Text   string    `json:"text,omitempty"`
}

type PostPutResp struct {
	PostID    uuid.UUID `json:"post_id,omitempty"`
	BlogID    uuid.UUID `json:"blog_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type PostDeleteReq struct {
	PostID uuid.UUID
	BlogID uuid.UUID
}
