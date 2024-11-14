package server

import (
	"net/http"
	"noteapp/internal/api"
	"noteapp/internal/repository"
	"noteapp/internal/service"
	"noteapp/pkg/database"
	"noteapp/pkg/logger"
	"strconv"
)

type server struct {
	router *http.ServeMux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func Start(c *configServer) {

	info := database.ConnectionInfo{
		Host: c.DataBase.Host,
		//Port:     c.Port,
		Username: c.DataBase.Username,
		DBName:   c.DataBase.Dbname,
		SSLMode:  c.DataBase.Sslmode,
	}
	db, err := database.NewPostgresConnection(info)
	if err != nil {
		logger.NewLog("server - Start()", 1, err, "Failed to create *sql.DB", info)
		return
	}
	defer db.Close()

	// init deps
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	noteRepo := repository.NewNotesRepository(db)

	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	noteService := service.NewNotesService(noteRepo)

	handler := api.NewHandler(userService, authService, noteService)

	srv := &server{
		router: handler.InitRouter(),
	}

	logger.NewLog("server - Start()", 5, nil, "Server starting...", nil)
	http.ListenAndServe(":"+strconv.Itoa(c.Port), srv)
}
