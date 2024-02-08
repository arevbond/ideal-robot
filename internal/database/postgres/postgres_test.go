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
			Port:     5430,
			Name:     "test_db",
		},
	}
	s, err := New(log, cfg)
	if err != nil {
		t.Fatalf("failed to create s: %v", err)
	}

	err = s.db.Ping()
	if err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	storage = s
	createTables(t)
}

// TODO: сделать автоматические миграции через goose
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

	_, err = storage.db.Exec(query)
	if err != nil {
		t.Fatal("Error executing SQL query:", err)
	}
}

func dropTables(t *testing.T) {
	q := `DROP TABLE hubs; DROP TABLE users;`

	_, err := storage.db.Exec(q)
	if err != nil {
		t.Fatal("Error executing SQL query:", err)
	}
}

func TestUser(t *testing.T) {
	// create
	users := []*models.User{
		{Username: "Nikita", PasswordHash: "OLDPASSWORDHASH", Email: "random@email.com", CreatedAt: time.Now()},
		{Username: "Nikita2", PasswordHash: "OLDPASSWORDHASH", Email: "random@em2ail.com", CreatedAt: time.Now()},
		{Username: "Nikita3", PasswordHash: "asldd2asdasd", Email: "random@em2ail.com", CreatedAt: time.Now()},
	}
	for _, u := range users {
		err := storage.CreateUser(context.Background(), u)
		if err != nil {
			t.Errorf("can't create user %s %s %s", u.Username, u.PasswordHash, u.Email)
		}
	}

	usersInDB, err := storage.GetUsers(context.Background())
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
		err = storage.UpdateUser(context.Background(), u)
		if err != nil {
			t.Errorf("can't update user %s %s %s", u.Username, u.PasswordHash, u.Email)
		}
	}

	// get by username
	fullUsers := []*models.DBUser{}
	for _, u := range newUsers {
		userInDB, err := storage.GetUserByUsername(context.Background(), u.Username)
		if err != nil {
			t.Error(fmt.Sprintf("can't get user %s from db by username", u.Username), err)
		}
		if u.PasswordHash != userInDB.PasswordHash {
			t.Errorf("wants: %s, got: %s", userInDB.PasswordHash, u.PasswordHash)
		}
		if u.Email != userInDB.Email {
			t.Errorf("wants: %s, got: %s", userInDB.Email, u.Email)
		}
		fullUsers = append(fullUsers, userInDB)
	}

	// get by id
	for _, u := range fullUsers {
		user, err := storage.GetUserByID(context.Background(), u.ID)
		if err != nil {
			t.Error(fmt.Sprintf("can't get user %s from db by id", u.Username), err)
		}
		if user.Username != u.Username {
			t.Error(fmt.Sprintf("incorrect getting user by id; want %s get %s", u.Username, user.Username), err)
		}
	}

	// delete
	for _, u := range newUsers[:1] {
		err = storage.DeleteUser(context.Background(), u)
		if err != nil {
			t.Error("can't delete user from db", err)
		}
	}

	finishUsers, err := storage.GetUsers(context.Background())
	if err != nil {
		t.Error("can't select users from db", err)
	}

	if len(finishUsers) != len(newUsers)-1 {
		t.Errorf("want users in db: %d, got: %d", len(newUsers)-1, len(finishUsers))
	}
}

func TestHub(t *testing.T) {
	user1 := &models.User{Username: "Nikita", PasswordHash: "SomeNikitaHash", Email: "nikita@mail.ru", CreatedAt: time.Now()}
	user2 := &models.User{Username: "Stas", PasswordHash: "SomeStasHash", Email: "stas@mail.ru", CreatedAt: time.Now()}
	user3 := &models.User{Username: "Akbar", PasswordHash: "SomeAkbarHash", Email: "akbar@mail.ru", CreatedAt: time.Now()}
	users := []*models.User{user1, user2, user3}

	dbUsers := []*models.DBUser{}
	hubs := []*models.Hub{}
	for i, u := range users {
		_ = storage.CreateUser(context.Background(), u)
		dbUser, _ := storage.GetUserByUsername(context.Background(), u.Username)
		dbUsers = append(dbUsers, dbUser)
		for j := 0; j < i+2; j++ {
			hubs = append(hubs, &models.Hub{OwnerID: dbUser.ID, Name: fmt.Sprintf("hub №%d", j),
				Description: fmt.Sprintf("hub by %s", dbUser.Username)})
		}
	}

	// create
	for _, h := range hubs {
		err := storage.CreateHub(context.Background(), h)
		if err != nil {
			t.Error("can't create hub", err)
		}
	}
	allHubsInDB, err := storage.GetHubs(context.Background())
	if err != nil {
		t.Error("can't get all hubs", err)
	}
	if len(allHubsInDB) != len(hubs) {
		t.Errorf("some hubs not create; want %d, get %d", len(hubs), len(allHubsInDB))
	}

	// get
	for _, h := range allHubsInDB {
		hInDb, err := storage.GetHubByID(context.Background(), h.ID)
		if err != nil {
			t.Error("can't get hub by id", err)
		}
		if h.Name != hInDb.Name || h.OwnerID != hInDb.OwnerID || h.Description != hInDb.Description {
			t.Errorf("incorrect hub in db; want: %s %s %s, get %s %s %s", h.Name, h.OwnerID, h.Description,
				hInDb.Name, hInDb.OwnerID, hInDb.Description)
		}
	}

	for i, u := range dbUsers {
		hubsByUserID, err := storage.GetHubsByUserID(context.Background(), u.ID)
		if err != nil {
			t.Error("can't get hubsByUserID by userID", err)
		}
		if len(hubsByUserID) != i+2 {
			t.Errorf("not enough hubsByUserID in user; want %d, get %d", i+2, len(hubsByUserID))
		}
	}

	for i, h := range allHubsInDB {
		if i < len(allHubsInDB)-1 {
			nextHub := allHubsInDB[i+1]
			newHub := &models.DBHub{ID: h.ID, OwnerID: nextHub.OwnerID,
				Name: nextHub.Name, Description: nextHub.Description}
			err := storage.UpdateHub(context.Background(), newHub)
			if err != nil {
				t.Error("can't update hub", err)
			}
		}
	}

	allHubsInDBold := allHubsInDB
	for i, h := range allHubsInDBold {
		if i < len(allHubsInDBold)-1 {
			nextHub := allHubsInDBold[i+1]
			curHubInDB, err := storage.GetHubByID(context.Background(), h.ID)
			if err != nil {
				t.Error("can't get hub by id", err)
			}
			if curHubInDB.Name != nextHub.Name {
				t.Errorf("don't update name hub; want %s get %s", nextHub.Name, curHubInDB.Name)
			}
			if curHubInDB.OwnerID != nextHub.OwnerID {
				t.Errorf("don't update owner_id hub; want %s get %s", nextHub.OwnerID, curHubInDB.OwnerID)
			}
			if curHubInDB.Description != nextHub.Description {
				t.Errorf("don't update descriprion hub; want %s get %s", nextHub.Description, curHubInDB.Description)
			}
		}
	}
	for _, h := range allHubsInDBold {
		err := storage.DeleteHub(context.Background(), h.ID)
		if err != nil {
			t.Error("can't delete hub by id", err)
		}
	}
	allHubs, err := storage.GetHubs(context.Background())
	if err != nil {
		t.Error("can't get all hubs", err)
	}
	if len(allHubs) != 0 {
		t.Errorf("not enough hubs in db; want %d, get %d", 0, len(allHubs))
	}
}
