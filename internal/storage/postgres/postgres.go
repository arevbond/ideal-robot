package postgres

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/lib/e"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) (*Storage, error) {
	dbSource := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Port, cfg.Storage.Name)
	conn, err := sqlx.Connect("pgx", dbSource)
	if err != nil {
		return nil, e.Wrap("connect to pgx failed", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, e.Wrap("ping failed", err)
	}
	return &Storage{db: conn, log: log}, nil
}
