package inmemory

import (
	"context"
	"sync"

	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/repo"
)

type Storage struct {
	mu    sync.RWMutex
	blogs map[string]model.Blog
	posts map[string]model.Post
}

func New() *Storage {
	return &Storage{
		blogs: make(map[string]model.Blog),
		posts: make(map[string]model.Post),
	}
}

func (s *Storage) GetBlog(_ context.Context, blogID string) (model.Blog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.blogs[blogID]; ok {
		return v, nil
	}

	return model.Blog{}, repo.ErrNotFound
}
func (s *Storage) AddBlog(_ context.Context, blog model.Blog) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blogs[blog.ID] = blog

	return blog.ID, nil
}
func (s *Storage) UpdateBlog(_ context.Context, blog model.Blog) (model.Blog, error) {
	s.mu.RLock()
	if _, ok := s.blogs[blog.ID]; !ok {
		s.mu.RUnlock()
		return model.Blog{}, repo.ErrNotFound
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blogs[blog.ID] = blog
	return blog, nil
}
func (s *Storage) DeleteBlog(_ context.Context, blogID string) error {
	s.mu.RLock()
	if _, ok := s.blogs[blogID]; !ok {
		s.mu.RUnlock()
		return repo.ErrNotFound
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.blogs, blogID)
	return nil
}
func (s *Storage) GetPost(_ context.Context, postID string) (model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.posts[postID]; ok {
		return v, nil
	}

	return model.Post{}, repo.ErrNotFound
}

func (s *Storage) GetPosts(_ context.Context, BlogID string) (blogs []model.Post, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.posts {
		if v.BlogID == BlogID {
			blogs = append(blogs, v)
		}
	}
	if len(blogs) == 0 {
		err = repo.ErrNotFound
		return
	}
	return
}

func (s *Storage) AddPost(_ context.Context, post model.Post) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts[post.ID] = post
	return post.ID, nil
}

func (s *Storage) UpdatePost(_ context.Context, post model.Post) (model.Post, error) {
	s.mu.RLock()
	if _, ok := s.posts[post.ID]; !ok {
		s.mu.RUnlock()
		return model.Post{}, repo.ErrNotFound
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts[post.ID] = post
	return model.Post{}, nil
}

func (s *Storage) DeletePost(_ context.Context, postID string) error {
	s.mu.RLock()
	if _, ok := s.posts[postID]; !ok {
		s.mu.RUnlock()
		return repo.ErrNotFound
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.posts, postID)
	return nil
}
