package model

import "time"

type User struct {
	ID string `json:"id,omitempty"`
}

type Blog struct {
	ID        string    `json:"id,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Post struct {
	ID        string    `json:"id,omitempty"`
	BlogID    string    `json:"blog_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// type Comment struct {
// 	ID        string    `json:"id,omitempty"`
// 	PostID    string    `json:"post_id,omitempty"`
// 	UserID    string    `json:"user_id,omitempty"`
// 	CreatedAt time.Time `json:"created_at,omitempty"`
// 	Text      string    `json:"text,omitempty"`
// }
