//nolint:all
package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rolan335/project/internal/apperror"
	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	defaultTtl  = time.Minute
	defaultSize = 100
)

func TestCache_GetBlog(t *testing.T) {
	existBlogID := uuid.New()
	existModel := model.DbBlog{
		ID:        existBlogID,
		UserID:    uuid.New(),
		Name:      gofakeit.Name(),
		CreatedAt: time.Now(),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockBlogRepository(ctrl)
	cache := NewCacheDecorator(defaultTtl, defaultSize, repository)

	repository.EXPECT().AddBlog(gomock.Any(), existModel).Return(existBlogID, nil).Times(1)
	cache.AddBlog(context.Background(), existModel)

	testCases := []struct {
		name            string
		in              uuid.UUID
		want            model.DbBlog
		err             error
		wantErr         error
		repositoryCalls int
	}{
		{
			name:            "found in cache",
			in:              existBlogID,
			want:            existModel,
			wantErr:         nil,
			repositoryCalls: 0,
		},
		{
			name:            "not found in cache",
			in:              uuid.New(),
			want:            model.DbBlog{},
			wantErr:         apperror.ErrNotFound,
			repositoryCalls: 1,
		},
		{
			name:            "nil uuid",
			in:              uuid.Nil,
			want:            model.DbBlog{},
			wantErr:         apperror.ErrNotFound,
			repositoryCalls: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repository.EXPECT().GetBlog(gomock.Any(), tt.in).Return(tt.want, tt.wantErr).Times(tt.repositoryCalls)
			got, err := cache.GetBlog(context.Background(), tt.in)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Cache.GetBlog() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Cache.GetBlog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_AddBlog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockBlogRepository(ctrl)
	cache := NewCacheDecorator(defaultTtl, defaultSize, repository)
	ids := map[string]uuid.UUID{"valid add": uuid.New()}
	testCases := []struct {
		name            string
		in              model.DbBlog
		want            uuid.UUID
		wantErr         error
		repositoryCalls int
	}{
		{
			name:            "valid add",
			in:              model.DbBlog{ID: ids["valid add"], UserID: uuid.New(), Name: gofakeit.Name(), CreatedAt: time.Now()},
			want:            ids["valid add"],
			repositoryCalls: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repository.EXPECT().AddBlog(gomock.Any(), tt.in).Return(tt.want, tt.wantErr).Times(tt.repositoryCalls)
			got, err := cache.AddBlog(context.Background(), tt.in)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Cache.AddBlog() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Cache.AddBlog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO
func TestCache_UpdateBlog(t *testing.T) {
	existBlogID := uuid.New()
	existModel := model.DbBlog{
		ID:        existBlogID,
		UserID:    uuid.New(),
		Name:      gofakeit.Name(),
		CreatedAt: time.Now(),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockBlogRepository(ctrl)
	cache := NewCacheDecorator(defaultTtl, defaultSize, repository)

	repository.EXPECT().AddBlog(gomock.Any(), existModel).Return(existBlogID, nil).Times(1)
	cache.AddBlog(context.Background(), existModel)

	testCases := []struct {
		name            string
		in              model.DbBlog
		want            model.DbBlog
		wantErr         error
		repositoryCalls int
	}{
		{
			name: "valid update",
			in:   existModel,
			want: model.DbBlog{
				ID:        existModel.ID,
				UserID:    existModel.UserID,
				Name:      "newName",
				CreatedAt: existModel.CreatedAt},
			repositoryCalls: 1,
		},
		{
			name: "not found",
			in: model.DbBlog{
				ID:        existBlogID,
				UserID:    uuid.New(),
				Name:      gofakeit.Name(),
				CreatedAt: time.Now(),
			},
			wantErr:         apperror.ErrNotFound,
			repositoryCalls: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repository.EXPECT().UpdateBlog(gomock.Any(), tt.in).Return(tt.want, tt.wantErr).Times(tt.repositoryCalls)
			got, err := cache.UpdateBlog(context.Background(), tt.in)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Cache.UpdateBlog() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Cache.UpdateBlog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_DeleteBlog(t *testing.T) {
	existBlogID := uuid.New()
	existModel := model.DbBlog{
		ID:        existBlogID,
		UserID:    uuid.New(),
		Name:      gofakeit.Name(),
		CreatedAt: time.Now(),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockBlogRepository(ctrl)
	cache := NewCacheDecorator(defaultTtl, defaultSize, repository)

	repository.EXPECT().AddBlog(gomock.Any(), existModel).Return(existBlogID, nil).Times(1)
	cache.AddBlog(context.Background(), existModel)

	testCases := []struct {
		name            string
		in              uuid.UUID
		wantErr         error
		repositoryCalls int
	}{
		{
			name:            "valid delete",
			in:              existBlogID,
			wantErr:         nil,
			repositoryCalls: 1,
		},
		{
			name:            "not found",
			in:              uuid.Nil,
			wantErr:         apperror.ErrNotFound,
			repositoryCalls: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repository.EXPECT().DeleteBlog(gomock.Any(), tt.in).Return(tt.wantErr).Times(tt.repositoryCalls)
			err := cache.DeleteBlog(context.Background(), tt.in)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Cache.DeleteBlog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_GoPollDeletion(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repository := mocks.NewMockBlogRepository(ctrl)
	customTTL := time.Nanosecond
	cache := NewCacheDecorator(customTTL, defaultSize, repository)
	addCount := 100
	// add 100 elements to cache
	for range addCount {
		blog := model.DbBlog{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Name:      gofakeit.Name(),
			CreatedAt: time.Now(),
		}
		repository.EXPECT().AddBlog(gomock.Any(), blog).Return(blog.ID, nil).Times(1)
		cache.AddBlog(context.Background(), blog)
	}
	_, len := cache.GetBlogLen()
	a.EqualValues(len, addCount)
	deleteInterval := time.Millisecond
	reallockInterval := time.Minute
	cache.GoPollDeletion(context.Background(), deleteInterval, reallockInterval)
	time.Sleep(time.Second)
	_, len = cache.GetBlogLen()
	a.EqualValues(0, len)
}
