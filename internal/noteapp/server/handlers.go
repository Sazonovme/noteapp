package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"noteapp/internal/noteapp/user"
	"noteapp/pkg/logger"
)

var (
	errMethodNotAllowed       = errors.New("method not allowed")
	errNoUserData             = errors.New("login or password not filled in")
	errInvalidPasswordOrLogin = errors.New("invalid password or login")
)

func (s *server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method != http.MethodPost {
		s.error(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		logger.NewLog("server", "handlerCreateUser", errMethodNotAllowed, nil, 5, "Method: "+r.Method+" not allowed")
		return
	}

	u := user.New()
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil || u.Login == "" || u.Password == "" {
		s.error(w, r, http.StatusNoContent, errNoUserData)
	}

	err = s.store.UserRepository.CreateUser(u)
	if err != nil {
		s.error(w, r, http.StatusConflict, err)
	}

}

func (s *server) handlerAuthUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method != http.MethodGet {
		s.error(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		logger.NewLog("server", "handlerCreateUser", errMethodNotAllowed, nil, 5, "Method: "+r.Method+" not allowed")
		return
	}
	u := user.New()
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
	}

	u2, err := s.store.UserRepository.FindByLogin(u.Login)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, errInvalidPasswordOrLogin)
	}

	err = u2.ComparePassword(u.Password)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, errInvalidPasswordOrLogin)
	}

}

// //////// Respond error //////////
func (s *server) error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	s.respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// //////// MiddleWare //////////
type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func middlewareCreateUser() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			f(w, r)
		}
	}
}
