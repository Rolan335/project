package main

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/Rolan335/project/config"
	"github.com/Rolan335/project/internal/app"
	"github.com/Rolan335/project/internal/cache"
	"github.com/Rolan335/project/internal/handler"
	"github.com/Rolan335/project/internal/repository"
	"github.com/Rolan335/project/internal/storage/pgconn"
	"github.com/Rolan335/project/internal/tracer"
	"github.com/Rolan335/project/internal/usecase"
	"github.com/Rolan335/project/migrations"
)

const configPath = "internal/config/config.yaml"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	if err := migrations.Migrate(cfg.Postgres.ConnStr); err != nil {
		log.Fatal().Err(err).Msg("")
	}
	conn, err := pgconn.GetConn(cfg.Postgres.ConnStr)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	defer conn.Close()

	tp, err := tracer.SetupTracer()
	if err != nil {
		log.Err(err).Msg("")
	}
	defer tp.Shutdown(context.Background())

	blogRepo := repository.NewBlogRepo(conn)

	//TODO: Перенести в конфиг
	cacheSize := 8
	deleteInterval := time.Second * 30
	ttl := time.Second * 15
	cache := cache.NewCacheDecorator(ttl, cacheSize, blogRepo)
	cache.GoPollDeletion(ctx, deleteInterval)

	blog := usecase.NewBlogProvider(cache)

	validate := validator.New()
	handle := handler.New(blog, validate)

	app := app.GetRouter(handle)

	if err := app.Listen(cfg.App.Port); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
