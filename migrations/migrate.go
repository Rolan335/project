package migrations

import (
	"database/sql"
	"embed"
	"log"

	// Import for side effects - needed for initializing the PostgresSQL driver
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

var sqlPath = "."

func Migrate(url string) error {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return errors.Wrap(err, "cannot connect to db")
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping db")
	}
	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "cannot set migrations dialect")
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return errors.Wrap(err, "cannot get migration version")
	}

	err = goose.Up(db, sqlPath)
	if err != nil {
		log.Println(err)
		if err = goose.DownTo(db, sqlPath, version); err != nil {
			return errors.Wrap(err, "cannot rollback migrations")
		}

		return errors.Wrap(err, "cannot up migrations")
	}

	return nil
}
