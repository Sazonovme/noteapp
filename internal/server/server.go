package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"noteapp/internal/api"
	"noteapp/internal/database"
	"noteapp/internal/repository"
	"noteapp/internal/service"
	"noteapp/pkg/logger"
	"os"
	"time"
)

func Start(ctx context.Context) error {
	data, err := os.ReadFile("../../configs/config.json")
	if err != nil {
		logger.NewLog("server - Start()", 1, nil, "Filed to read config.json", nil)
		return err
	}

	config := NewConfig()
	if err = json.Unmarshal(data, config); err != nil {
		logger.NewLog("server - Start()", 1, nil, "Filed to unmarshal JSON", nil)
		return err
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
		return err
	}
	defer db.Close()

	// init deps
	userRepo := repository.NewUserRepository(db)
	noteRepo := repository.NewNotesRepository(db)

	userService := service.NewUserService(userRepo)
	noteService := service.NewNotesService(noteRepo)

	handler := api.NewHandler(userService, noteService)

	srv := &http.Server{
		Addr:    config.Addr,
		Handler: handler.InitHandler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.NewLog("server - Start()", 1, err, "Filed to start server", info)
		}
	}()

	logger.NewLog("server - Start()", 5, nil, "Server starting...", nil)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
