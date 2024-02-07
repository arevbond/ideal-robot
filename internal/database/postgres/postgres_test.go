package postgres

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/models"
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"
)

var storage *Storage

func TestMain(m *testing.M) {
	t := &testing.T{}
	setup(t)
	code := m.Run()
	teardown(t)
	os.Exit(code)
}

func teardown(t *testing.T) {
	dropTables(t)
}

func setup(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := &config.Config{
		DB: config.Database{
			Username: "admin",
			Password: "admin123",
			Name:     "test_db",
		},
	}
	s, err := New(log, cfg)
	if err != nil {
		t.Fatalf("failed to create s: %v", err)
	}

	// Проверка соединения с базой данных
	err = s.db.Ping()
	if err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	storage = s
	createTables(t)
}

func createTables(t *testing.T) {
	sqlFile, err := os.Open("../../../migrations/20240205101902_init.sql")
	if err != nil {
		t.Fatal("Error opening SQL file:", err)
	}
	defer sqlFile.Close()
	queryBytes, err := io.ReadAll(sqlFile)
	if err != nil {
		t.Fatal("Error reading SQL file:", err)
	}
	query := string(queryBytes)

	// Выполняем SQL-запрос
	_, err = storage.db.Exec(query)
	if err != nil {
		t.Fatal("Error executing SQL query:", err)
	}
}

func dropTables(t *testing.T) {
	q := `DROP TABLE hubs; DROP TABLE users;`

	// Выполняем SQL-запрос
	_, err := storage.db.Exec(q)
	if err != nil {
		t.Fatal("Error executing SQL query:", err)
	}
}

func TestUser(t *testing.T) {
	// create
	users := []*models.DBUser{
		{Username: "Nikita", PasswordHash: "OLDPASSWORDHASH", Email: "random@email.com", CreatedAt: time.Now()},
		{Username: "Nikita2", PasswordHash: "OLDPASSWORDHASH", Email: "random@em2ail.com", CreatedAt: time.Now()},
		{Username: "Nikita3", PasswordHash: "asldd2asdasd", Email: "random@em2ail.com", CreatedAt: time.Now()},
	}
	for _, u := range users {
		//fmt.Println(u)
		err := storage.CreateUser(context.Background(), u)
		if err != nil {
			t.Errorf("can't create user %s %s %s", u.Username, u.PasswordHash, u.Email)
		}
	}

	q := `SELECT * FROM users`
	usersInDB := []*models.DBUser{}
	err := storage.db.Select(&usersInDB, q)
	if err != nil {
		t.Error("can't select users from db", err)
	}

	if len(usersInDB) != len(users) {
		t.Error("not enough created users", err)
	}

	// update
	newUsers := []*models.DBUser{
		{Username: "Nikita", PasswordHash: "NEWPASSWORDHASH", Email: "randomNEW@email.com", CreatedAt: time.Now()},
		{Username: "Nikita2", PasswordHash: "NEWPASSWORDHASH", Email: "randomBEW@em2ail.com", CreatedAt: time.Now()},
		{Username: "Nikita3", PasswordHash: "123", Email: "randomBEW@em2ail.com", CreatedAt: time.Now()},
	}
	for _, u := range newUsers {
		//fmt.Println(u)
		err = storage.UpdateUser(context.Background(), u)
		if err != nil {
			t.Errorf("can't update user %s %s %s", u.Username, u.PasswordHash, u.Email)
		}
	}

	for _, u := range newUsers {
		userInDB, err := storage.GetUser(context.Background(), u.Username)
		if err != nil {
			t.Error(fmt.Sprintf("can't get user %s from db", u.Username), err)
		}
		if u.PasswordHash != userInDB.PasswordHash {
			t.Errorf("wants: %s, got: %s", userInDB.PasswordHash, u.PasswordHash)
		}
		if u.Email != userInDB.Email {
			t.Errorf("wants: %s, got: %s", userInDB.Email, u.Email)
		}
	}

	for _, u := range newUsers[:1] {
		err = storage.DeleteUser(context.Background(), u)
		if err != nil {
			t.Error("can't delete user from db", err)
		}
	}

	finishUsers := []*models.DBUser{}
	err = storage.db.Select(&finishUsers, q)
	if err != nil {
		t.Error("can't select users from db", err)
	}

	if len(finishUsers) != len(newUsers)-1 {
		t.Errorf("want users in db: %d, got: %d", len(newUsers)-1, len(finishUsers))
	}

}
