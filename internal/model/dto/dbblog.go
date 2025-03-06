package dto

import (
	"time"

	"github.com/google/uuid"
)

type DbUser struct {
	ID uuid.UUID `json:"id,omitempty" db:"id"`
}

type DbBlog struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	UserID    uuid.UUID `json:"user_id,omitempty" db:"users_id"`
	Name      string    `json:"name,omitempty" db:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type DbPost struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	BlogID    uuid.UUID `json:"blog_id,omitempty" db:"blogs_id"`
	Title     string    `json:"title,omitempty" db:"title"`
	Text      string    `json:"text,omitempty" db:"text"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
