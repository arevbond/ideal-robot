package postgres

import (
	"HestiaHome/internal/database"
	"HestiaHome/internal/lib/e"
	"HestiaHome/internal/models"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log/slog"
)

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
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

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.DBUser, error) {
	q := `SELECT * FROM users WHERE username = $1`
	var user models.DBUser

	err := s.db.GetContext(ctx, &user, q, username)

	if err == sql.ErrNoRows {
		return nil, database.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap("can't get user from storage by username", err)
	}

	return &user, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (*models.DBUser, error) {
	q := `SELECT * FROM users WHERE id = $1`

	var user models.DBUser

	err := s.db.GetContext(ctx, &user, q, id)

	if err == sql.ErrNoRows {
		return nil, database.ErrUserNotExist
	}

	if err != nil {
		return nil, e.Wrap("can't get user from storage by id", err)
	}

	return &user, nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]*models.DBUser, error) {
	q := `SELECT * FROM users`

	users := []*models.DBUser{}

	err := s.db.SelectContext(ctx, &users, q)
	if err != nil {
		return nil, e.Wrap("can't get all users from storage", err)
	}

	return users, nil
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
