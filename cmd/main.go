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
	"github.com/Rolan335/project/internal/metric"
	"github.com/Rolan335/project/internal/repository"
	"github.com/Rolan335/project/internal/storage/pgconn"
	"github.com/Rolan335/project/internal/tracer"
	"github.com/Rolan335/project/internal/usecase"
	"github.com/Rolan335/project/migrations"
)

const configPath = "config/config.yaml"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	if err := migrations.Migrate(cfg.Postgres.ConnStr); err != nil {
		log.Panic().Err(err).Msg("")
	}
	conn, err := pgconn.GetConn(cfg.Postgres.ConnStr)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}
	defer conn.Close()

	tp, err := tracer.SetupTracer()
	if err != nil {
		log.Err(err).Msg("")
	}
	defer tp.Shutdown(context.Background()) //nolint

	blogRepo := repository.NewBlogRepo(conn)

	//TODO: Перенести в конфиг
	cacheSize := 100
	ttl := time.Second * 15
	cache := cache.NewCacheDecorator(ttl, cacheSize, blogRepo)
	deleteInterval := time.Second * 30
	realocInterval := time.Minute
	cache.GoPollDeletion(ctx, deleteInterval, realocInterval)

	blog := usecase.NewBlogProvider(cache)

	metric.MustRegisterMetrics()
	pollInterval := 10 * time.Second
	metric.GoCountCacheLen(ctx, pollInterval, cache)

	validate := validator.New()
	handle := handler.New(blog, validate)

	apiEndpoint := app.GetRouter(handle)

	metricEndpoint := app.GetMetricsRouter()

	go func() {
		if err := metricEndpoint.Listen(":8081"); err != nil {
			log.Panic().Err(err).Msg("")
		}
	}()

	if err := apiEndpoint.Listen(cfg.App.Port); err != nil {
		log.Panic().Err(err).Msg("")
	}
}
