package cachedata

import (
	"time"

	"github.com/Rolan335/project/internal/model"
)

type CacheBlog struct {
	Deadline time.Time
	Db  model.DbBlog
}

type CachePost struct {
	Deadline time.Time
	Db  model.DbPost
}
