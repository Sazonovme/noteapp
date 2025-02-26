package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"noteapp/internal/database"
	"noteapp/internal/repository"
	"noteapp/internal/service"
	"os"
	"testing"
)

func NewTestHandler(t *testing.T) (
	*http.ServeMux,
	func(*sql.DB, *testing.T, string),
	func(*sql.DB, *testing.T, string),
	*sql.DB,
) {
	t.Helper()

	data, err := os.ReadFile("../../configs/config.json")
	if err != nil {
		t.Fatal(err)
	}

	config := struct {
		LogLevel int    `json:"log-level"`
		Port     int    `json:"port"`
		Addr     string `json:"addr"`
		DataBase struct {
			Host     string `json:"host"`
			Username string `username:"username"`
			Dbname   string `json:"dbname"`
			Sslmode  string `json:"sslmode"`
		} `json:"database"`
	}{}
	if err = json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}

	info := database.ConnectionInfo{
		Host: config.DataBase.Host,
		//Port:     c.Port,
		Username: config.DataBase.Username,
		DBName:   config.DataBase.Dbname,
		SSLMode:  config.DataBase.Sslmode,
	}
	db, create, teardown := database.NewTestPostgresConnection(t, info)

	// init deps
	//authRepo := repository.NewTestAuthRepository(db)
	userRepo := repository.NewTestUserRepository(db)
	noteRepo := repository.NewTestNotesRepository(db)

	//authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	noteService := service.NewNotesService(noteRepo)

	h := NewHandler(userService, noteService)
	handler := h.InitHandler()
	return handler, create, teardown, db
}
