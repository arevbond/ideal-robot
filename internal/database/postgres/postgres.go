package postgres

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/database"
	"HestiaHome/internal/lib/e"
	"HestiaHome/internal/models"
	"context"
	"database/sql"
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
	dbSource := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Name)
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

func (s *Storage) CreateUser(ctx context.Context, user *models.DBUser) error {
	q := `INSERT INTO users (username, password_hash, email, created_at) VALUES ($1, $2, $3, $4)`
	if _, err := s.db.ExecContext(ctx, q, user.Username, user.PasswordHash, user.Email, user.CreatedAt); err != nil {
		return e.Wrap("can't create user in database", err)
	}
	s.log.Debug("create user", slog.String("username", user.Username))
	return nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.DBUser) error {
	q := `UPDATE users SET username = $1, password_hash = $2, email = $3 WHERE username = $4`
	_, err := s.db.ExecContext(ctx, q, user.Username, user.PasswordHash, user.Email, user.Username)

	if err != nil {
		return e.Wrap("can't update user", err)
	}
	s.log.Debug("update user", slog.String("username", user.Username))
	return nil
}

func (s *Storage) GetUser(ctx context.Context, username string) (*models.DBUser, error) {
	q := `SELECT * FROM users WHERE username = $1`
	var user models.DBUser

	err := s.db.GetContext(ctx, &user, q, username)

	if err == sql.ErrNoRows {
		return nil, database.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap("can't get user from storage", err)
	}

	return &user, nil
}

func (s *Storage) DeleteUser(ctx context.Context, user *models.DBUser) error {
	q := `DELETE FROM users WHERE username = $1`

	_, err := s.db.ExecContext(ctx, q, user.Username)
	if err != nil {
		return e.Wrap("can't delete user", err)
	}
	s.log.Debug("delete user", slog.String("username", user.Username))
	return nil
}

func (s *Storage) CreateHub(ctx context.Context, hub *models.Hub) error {
	q := `INSERT INTO hubs (user_id, name, description) VALUE ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, q, hub.OwnerID, hub.Name, hub.Description)

	if err != nil {
		return e.Wrap("can't create hub", err)
	}
	return nil
}
