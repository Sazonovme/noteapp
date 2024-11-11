package server

import (
	"net/http"
	"noteapp/internal/noteapp/sqlstore"
	"noteapp/pkg/logger"
)

type server struct {
	store   *sqlstore.Sqlstore
	handler *http.ServeMux
}

func NewServer(store *sqlstore.Sqlstore) *server {
	return &server{
		store: store,
	}
}

func (s *server) configureHandler() {
	srv := http.NewServeMux()
	srv.HandleFunc("/sign-up", ChainMiddleware(
		s.handlerCreateUser,
		s.middlewareNoCors()))
	srv.HandleFunc("/sign-in", ChainMiddleware(
		s.handlerAuthUser,
		s.middlewareNoCors()))
	srv.HandleFunc("/notes", ChainMiddleware(
		s.handlerNotes,
		s.middlewareNoCors(),
		s.middlewareAuth()))
	srv.HandleFunc("/refresh-token",
		s.handlerRefreshToken,
	)
	s.handler = srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func Start(c *ConfigServer) {

	store, err := sqlstore.NewStore(c.Db_connection)
	if err != nil {
		logger.NewLog("server.go - Start() - sqlstore.NewStore()", 1, err, "Failed to create *sqlstore.Sqlstore", nil)
		return
	}
	defer store.Db.Close()

	srv := NewServer(store)

	srv.configureHandler()

	srv.store.UserRepository = sqlstore.NewUserRepository(srv.store)

	logger.NewLog("server.go - Start() - http.ListenAndServe()", 5, nil, "Server starting...", nil)

	err = http.ListenAndServe(":"+c.Port, srv)
	if err != nil {
		logger.NewLog("server.go - Start() - http.ListenAndServe()", 1, err, "Failed to start server", nil)
		return
	}
}
