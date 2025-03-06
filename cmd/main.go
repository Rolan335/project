package main

import (
	"github.com/Rolan335/project/internal/app"
	"github.com/Rolan335/project/internal/config"
	"github.com/Rolan335/project/internal/handler"
	"github.com/Rolan335/project/internal/repository/blogprovider"
	"github.com/Rolan335/project/internal/storage/pgconn"
	"github.com/Rolan335/project/internal/usecase"
	"github.com/Rolan335/project/migrations"
)

func main() {
	//TODO: перенести в env
	cfg, err := config.New()
	if err != nil {
		panic("main.config: " + err.Error())
	}

	if err := migrations.Migrate(cfg.Postgres.ConnStr); err != nil {
		panic("main.migrations.Migrate: " + err.Error())
	}
	conn, err := pgconn.GetConn(cfg.Postgres.ConnStr)
	if err != nil {
		panic("main.pgconn.New: " + err.Error())
	}
	repository := blogprovider.New(conn)

	blog := usecase.NewBlogProvider(repository)

	handle := handler.New(blog)

	app := app.GetRouter(handle)

	if err := app.Listen(cfg.App.Port); err != nil {
		panic("main.app.Listen: " + err.Error())
	}
}
