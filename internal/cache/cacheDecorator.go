package cache

import (
	"context"
	"sync"
	"time"

	"github.com/Rolan335/project/internal/cache/cachedata"
	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type CacheDecorator struct {
	size       int
	blogCache  *cachedata.BlogCache
	postCache  *cachedata.PostCache
	repository repository.BlogRepository
	once       sync.Once
}

func NewCacheDecorator(ttl time.Duration, size int, repository repository.BlogRepository) *CacheDecorator {
	return &CacheDecorator{
		size:       size,
		blogCache:  cachedata.NewBlogCache(size, ttl),
		postCache:  cachedata.NewPostCache(size, ttl),
		repository: repository,
	}
}

// Returns cache name and len
func (c *CacheDecorator) GetBlogLen() (string, int) {
	return "blogCache", c.blogCache.GetLen()
}

// Returns cache name and len
func (c *CacheDecorator) GetPostLen() (string, int) {
	return "postCache", c.postCache.GetLen()
}

func (c *CacheDecorator) GoPollDeletion(ctx context.Context, deleteInterval time.Duration, reallocInterval time.Duration) {
	c.once.Do(func() {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(deleteInterval):
					c.blogCache.DeleteExpired()
					c.postCache.DeleteExpired()
				// мапа в го не уменьшается по размеру при удалении ключей, чтобы не умереть по памяти, реаллоцируем
				case <-time.Tick(reallocInterval):
					if _, len := c.GetBlogLen(); len > c.size {
						c.blogCache.DeleteFull()
					}
					if _, len := c.GetPostLen(); len > c.size {
						c.postCache.DeleteFull()
					}
				}
			}
		}()
	})
}

func (c *CacheDecorator) GetBlog(ctx context.Context, blogID uuid.UUID) (model.DbBlog, error) {
	tracer := otel.Tracer("project")
	ctx, span := tracer.Start(ctx, "Cache")
	defer span.End()
	if model, ok := c.blogCache.Get(ctx, blogID); ok {
		log.Debug().Str("uuid:", blogID.String()).Msg("cache hit")
		return model, nil
	}
	log.Debug().Str("uuid:", blogID.String()).Msg("cache miss")
	blog, err := c.repository.GetBlog(ctx, blogID)
	if err != nil {
		return model.DbBlog{}, errors.Wrap(err, "cacheDecorator.GetBlog")
	}
	// add to cache if in db, but not in cache
	c.blogCache.Set(ctx, blog)
	return blog, nil
}
func (c *CacheDecorator) AddBlog(ctx context.Context, blog model.DbBlog) (uuid.UUID, error) {
	id, err := c.repository.AddBlog(ctx, blog)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cacheDecorator.AddBlog")
	}
	// set to cache only if success insert into repo
	c.blogCache.Set(ctx, blog)
	return id, nil
}
func (c *CacheDecorator) UpdateBlog(ctx context.Context, blog model.DbBlog) (model.DbBlog, error) {
	newBlog, err := c.repository.UpdateBlog(ctx, blog)
	if err != nil {
		return model.DbBlog{}, errors.Wrap(err, "cacheDecorator.UpdateBlog")
	}
	// update cache only if success into repo
	c.blogCache.Set(ctx, blog)
	return newBlog, nil
}

// TODO: При удалении блога остаются посты блога в кэше cachedata.PostCache
func (c *CacheDecorator) DeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	err := c.repository.DeleteBlog(ctx, blogID)
	if err != nil {
		return errors.Wrap(err, "cacheDecorator.DeleteBlog")
	}
	c.blogCache.Delete(ctx, blogID)
	return nil
}

func (c *CacheDecorator) GetPost(ctx context.Context, postID uuid.UUID) (model.DbPost, error) {
	if model, ok := c.postCache.Get(ctx, postID); ok {
		log.Debug().Str("uuid:", postID.String()).Msg("cache hit")
		return model, nil
	}
	log.Debug().Str("uuid:", postID.String()).Msg("cache miss")
	post, err := c.repository.GetPost(ctx, postID)
	if err != nil {
		return model.DbPost{}, errors.Wrap(err, "cacheDecorator.GetPost")
	}
	// add to cache if in db, but not in cache
	c.postCache.Set(ctx, post)
	return post, nil
}

func (c *CacheDecorator) GetPosts(ctx context.Context, blogID uuid.UUID) ([]model.DbPost, error) {
	// не идём в кэш, так как там могут быть не все посты и в любом случае обращение в бд.
	return c.repository.GetPosts(ctx, blogID)
}
func (c *CacheDecorator) AddPost(ctx context.Context, post model.DbPost) (uuid.UUID, error) {
	id, err := c.repository.AddPost(ctx, post)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "cacheDecorator.AddPost")
	}
	// set to cache only if success insert into repo
	c.postCache.Set(ctx, post)
	return id, nil
}
func (c *CacheDecorator) UpdatePost(ctx context.Context, post model.DbPost) (model.DbPost, error) {
	newPost, err := c.repository.UpdatePost(ctx, post)
	if err != nil {
		return model.DbPost{}, errors.Wrap(err, "cacheDecorator.UpdatePost")
	}
	// update cache only if success into repo
	c.postCache.Set(ctx, post)
	return newPost, nil
}
func (c *CacheDecorator) DeletePost(ctx context.Context, postID uuid.UUID, blogID uuid.UUID) error {
	err := c.repository.DeletePost(ctx, postID, blogID)
	if err != nil {
		return errors.Wrap(err, "cacheDecorator.DeletePost")
	}
	c.postCache.Delete(ctx, postID)
	return nil
}
