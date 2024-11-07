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
		s.middlewareNoCors(),
		s.middlewareAddheadersNoCors()))
	srv.HandleFunc("/sign-in", ChainMiddleware(
		s.handlerAuthUser,
		s.middlewareNoCors(),
		s.middlewareAddheadersNoCors()))
	srv.HandleFunc("/notes", ChainMiddleware(
		s.handlerNotes,
		s.middlewareNoCors(),
		//s.middlewareAddheadersNoCors()))
		s.middlewareAuth()))
	s.handler = srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func Start(c *ConfigServer) error {

	// connect store
	store, err := sqlstore.NewStore(c.Db_connection)
	if err != nil {
		logger.NewLog("server", "Start", err, nil, 1, "Server can`t start")
	}

	defer store.Db.Close()

	// new server
	srv := NewServer(store)

	//configure router server
	srv.configureHandler()

	// configure user repository
	srv.store.UserRepository = sqlstore.NewUserRepository(srv.store)

	// start server
	logger.NewLog("server", "Start", nil, nil, 6, "Server starting...")
	err = http.ListenAndServe(c.Port, srv)
	if err != nil {
		return err
	}
	logger.NewLog("server", "Start", nil, nil, 6, "Server stoped")
	return nil
}
