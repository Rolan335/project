package cachedata

import (
	"context"
	"sync"
	"time"

	"github.com/Rolan335/project/internal/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BlogCache struct {
	ttl  time.Duration
	size int
	mu   *sync.RWMutex
	data map[uuid.UUID]CacheBlog
}

func NewBlogCache(size int, ttl time.Duration) *BlogCache {
	return &BlogCache{
		ttl:  ttl,
		size: size,
		mu:   &sync.RWMutex{},
		data: make(map[uuid.UUID]CacheBlog, size), /* prealloc memory */
	}
}

func (b *BlogCache) GetLen() int {
	return len(b.data)
}

func (b *BlogCache) Get(_ context.Context, uuid uuid.UUID) (model.DbBlog, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if v, ok := b.data[uuid]; ok {
		return v.Db, true
	}
	return model.DbBlog{}, false
}

func (b *BlogCache) Set(_ context.Context, model model.DbBlog) uuid.UUID {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[model.ID] = CacheBlog{
		Deadline: time.Now().Add(b.ttl),
		Db:       model,
	}
	return model.ID
}

func (b *BlogCache) Delete(_ context.Context, uuid uuid.UUID) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.data[uuid]; ok {
		delete(b.data, uuid)
		return true
	}
	return false
}

func (b *BlogCache) DeleteExpired() {
	b.mu.Lock()
	defer b.mu.Unlock()
	start := time.Now()
	for k, v := range b.data {
		if start.After(v.Deadline) {
			log.Debug().Str("uuid:", k.String()).Msg("deleted by deadline")
			delete(b.data, k)
		}
	}
}

func (b *BlogCache) DeleteFull() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = make(map[uuid.UUID]CacheBlog, b.size)
}

type PostCache struct {
	size int
	ttl  time.Duration
	mu   *sync.RWMutex
	data map[uuid.UUID]CachePost
}

func NewPostCache(size int, ttl time.Duration) *PostCache {
	return &PostCache{
		size: size,
		ttl:  ttl,
		mu:   &sync.RWMutex{},
		data: make(map[uuid.UUID]CachePost, size), /* prealloc memory */
	}
}

func (b *PostCache) GetLen() int {
	return len(b.data)
}

func (b *PostCache) Get(_ context.Context, uuid uuid.UUID) (model.DbPost, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if v, ok := b.data[uuid]; ok {
		return v.Db, true
	}
	return model.DbPost{}, false
}

func (b *PostCache) Set(_ context.Context, model model.DbPost) uuid.UUID {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[model.ID] = CachePost{
		Deadline: time.Now().Add(b.ttl),
		Db:       model,
	}
	return model.ID
}

func (b *PostCache) Delete(_ context.Context, uuid uuid.UUID) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.data[uuid]; ok {
		delete(b.data, uuid)
		return true
	}
	return false
}

func (b *PostCache) DeleteExpired() {
	b.mu.Lock()
	defer b.mu.Unlock()
	start := time.Now()
	for k, v := range b.data {
		if start.After(v.Deadline) {
			log.Debug().Str("uuid:", k.String()).Msg("deleted by deadline")
			delete(b.data, k)
		}
	}
}

func (b *PostCache) DeleteFull() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = make(map[uuid.UUID]CachePost, b.size)
}
