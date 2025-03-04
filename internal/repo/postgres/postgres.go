package postgres

import (
	"context"
	"fmt"

	"github.com/Rolan335/project/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_NAME"`
}

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg *Config) (*Storage, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}
	return &Storage{
		db: conn,
	}, nil
}

func (s *Storage) GetBlog(ctx context.Context, blogID string) (model.Blog, error) {
	return model.Blog{}, nil
}
func (s *Storage) AddBlog(ctx context.Context, blog model.Blog) (string, error) {
	return "", nil
}
func (s *Storage) UpdateBlog(ctx context.Context, blog model.Blog) (model.Blog, error) {
	return model.Blog{}, nil
}
func (s *Storage) DeleteBlog(ctx context.Context, blogID string) error {
	return nil
}
func (s *Storage) GetPost(ctx context.Context, postID string) (model.Post, error) {
	return model.Post{}, nil
}
func (s *Storage) GetPosts(ctx context.Context, BlogID string) ([]model.Post, error) {
	return nil, nil
}
func (s *Storage) AddPost(ctx context.Context, post model.Post) (string, error) {
	return "", nil
}
func (s *Storage) UpdatePost(ctx context.Context, post model.Post) (model.Post, error) {
	return model.Post{}, nil
}
func (s *Storage) DeletePost(ctx context.Context, postID string) error {
	return nil
}
