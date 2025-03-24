// nolint
package integration

import (
	"context"
	"os"
	"testing"

	"github.com/Rolan335/project/internal/apperror"
	"github.com/Rolan335/project/internal/model"
	"github.com/Rolan335/project/internal/repository"
	"github.com/Rolan335/project/internal/storage/pgconn"
	"github.com/Rolan335/project/internal/usecase"
	"github.com/Rolan335/project/migrations"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestBlogProvider(t *testing.T) {
	a := assert.New(t)
	godotenv.Load()

	pgConnStr := os.Getenv("POSTGRES_CONNSTR")

	if err := migrations.Migrate(pgConnStr); err != nil {
		log.Panic().Err(err).Msg("")
	}

	pg, err := pgconn.GetConn(pgConnStr)
	a.NoError(err)

	repository := repository.NewBlogRepo(pg)

	blogprovider := usecase.NewBlogProvider(repository)

	ctx := context.Background()

	addBlogReq := model.BlogPostReq{UserID: uuid.New(), Name: gofakeit.Name()}
	var addBlogResp model.BlogPostResp
	t.Run("AddBlog", func(t *testing.T) {
		var err error
		addBlogResp, err = blogprovider.AddBlog(ctx, addBlogReq)
		a.NoError(err)
		a.IsType(model.BlogPostResp{}, addBlogResp)
		//returns valid uuid
		a.NotZero(addBlogResp.BlogID.String())
	})

	getBlogReq := model.BlogGetReq{BlogID: addBlogResp.BlogID}
	var getBlogResp model.BlogGetResp
	t.Run("GetBlog", func(t *testing.T) {
		var err error
		getBlogResp, err = blogprovider.GetBlog(ctx, getBlogReq)
		a.NoError(err)
		a.IsType(model.BlogGetResp{}, getBlogResp)
		a.NotZero(getBlogResp.BlogID.String())
		a.NotZero(getBlogResp.UserID.String())

		resp, err := blogprovider.GetBlog(ctx, model.BlogGetReq{BlogID: uuid.Nil})
		a.Zero(resp)
		a.ErrorIs(err, apperror.ErrNotFound)
	})

	updateBlogReq := model.BlogPutReq{BlogID: addBlogResp.BlogID, UserID: addBlogReq.UserID, Name: gofakeit.Name()}
	t.Run("UpdateBlog", func(t *testing.T) {
		UpdateBlogResp, err := blogprovider.UpdateBlog(ctx, updateBlogReq)
		a.NoError(err)
		a.Equal(updateBlogReq.BlogID, UpdateBlogResp.BlogID)
		a.Equal(updateBlogReq.UserID, UpdateBlogResp.UserID)
		a.Equal(updateBlogReq.Name, UpdateBlogResp.Name)

		resp, err := blogprovider.GetBlog(ctx, model.BlogGetReq{BlogID: updateBlogReq.BlogID})
		a.NoError(err)
		a.Equal(UpdateBlogResp.BlogID, resp.BlogID)
		a.Equal(UpdateBlogResp.UserID, resp.UserID)
		a.Equal(UpdateBlogResp.Name, resp.Name)
		a.Equal(UpdateBlogResp.CreatedAt, resp.CreatedAt)
	})

	DeleteBlogReq := model.BlogDeleteReq{BlogID: addBlogResp.BlogID}
	t.Run("DeleteBlog", func(t *testing.T) {
		err := blogprovider.DeleteBlog(ctx, DeleteBlogReq)
		a.NoError(err)
		resp, err := blogprovider.GetBlog(ctx, model.BlogGetReq{BlogID: DeleteBlogReq.BlogID})
		a.ErrorIs(err, apperror.ErrNotFound)
		a.Zero(resp)

		err = blogprovider.DeleteBlog(ctx, model.BlogDeleteReq{BlogID: uuid.New()})
		a.ErrorIs(err, apperror.ErrNotFound)
	})
}
