package api

import (
	"encoding/json"
	"net/http"
	"noteapp/internal/repository"
	"noteapp/internal/service"
	"noteapp/pkg/database"
	"noteapp/pkg/logger"
	"os"
)

func NewTestHandler() (*http.ServeMux, error) {
	data, err := os.ReadFile("../../configs/config.json")
	if err != nil {
		logger.NewLog("server - Start()", 1, nil, "Filed to read config.json", nil)
		return nil, err
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
	if err = json.Unmarshal(data, config); err != nil {
		logger.NewLog("server - Start()", 1, nil, "Filed to unmarshal JSON", nil)
		return nil, err
	}

	info := database.ConnectionInfo{
		Host: config.DataBase.Host,
		//Port:     c.Port,
		Username: config.DataBase.Username,
		DBName:   config.DataBase.Dbname,
		SSLMode:  config.DataBase.Sslmode,
	}
	db, err := database.NewPostgresConnection(info)
	if err != nil {
		logger.NewLog("server - Start()", 1, err, "Failed to create *sql.DB", info)
		return nil, err
	}
	defer db.Close()
	// init deps
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	noteRepo := repository.NewNotesRepository(db)

	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	noteService := service.NewNotesService(noteRepo)

	h := NewHandler(userService, authService, noteService)
	handler := h.InitHandler()
	return handler, nil
}
