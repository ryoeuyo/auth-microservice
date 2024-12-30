package postgres

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	sql "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/ryoeuyo/sso/internal/config"
)

func MustInit(cfg *config.Database) *Database {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password,
	)

	pgxCfg, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Fatalf("failed parse config: %s", err.Error())
	}

	db := sql.OpenDB(*pgxCfg)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed ping database: %s", err.Error())
	}

	// run migrations
	if err := goose.SetDialect(cfg.Engine); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, cfg.MigrationDir); err != nil {
		log.Fatal(err)
	}

	return New(db)
}
