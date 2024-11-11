package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"noteapp/internal/noteapp/user"

	"github.com/asaskevich/govalidator"
)

var (
	errMethodNotAllowed = errors.New("method not allowed")
)

type requestData struct {
	Data struct {
		Login       string `json:"login" valid:"required"`
		Password    string `json:"password" valid:"required"`
		Fingerprint string `json:"fingerprint" valid:"required"`
	} `json:"data"`
}

func (s *server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		s.error(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		//logger.NewLog("server", "handlerCreateUser", errMethodNotAllowed, nil, 5, "Method: "+r.Method+" not allowed")
		return
	}

	d := &requestData{}
	err := json.NewDecoder(r.Body).Decode(d)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := govalidator.ValidateStruct(d); err != nil {
		s.respond(w, r, http.StatusBadRequest, err)
		return
	}

	u := user.New()
	u.Login = d.Data.Login
	u.Fingerprint = d.Data.Fingerprint
	u.Password = d.Data.Password
	err = s.store.UserRepository.CreateUser(u)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (s *server) handlerAuthUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		s.error(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		//logger.NewLog("server", "handlerCreateUser", errMethodNotAllowed, nil, 5, "Method: "+r.Method+" not allowed")
		return
	}

	d := &requestData{}
	err := json.NewDecoder(r.Body).Decode(d)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	u := user.New()
	u.Login = d.Data.Login
	u.Fingerprint = d.Data.Fingerprint
	u.Password = d.Data.Password

	u2, err := s.store.UserRepository.FindByLogin(u.Login)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	err = u2.ComparePassword(u.Password)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	refSession, err := s.store.AuthRepository.MakeRefreshSession(u.Login, u.Fingerprint)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// cookie := &http.Cookie{
	// 	HttpOnly: true,
	// 	Name:     "refreshToken",
	// 	Value:    refSession.RefreshToken,
	// }
	// http.SetCookie(w, cookie)

	m := map[string]string{
		"accessToken":  refSession.AccessToken,
		"refreshToken": refSession.RefreshToken,
		"exp":          refSession.Exp,
	}

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		//logger.NewLog("server", "handlerCreateUser", err, m, 6, "http response: json encode error")
		return
	}
}

func (s *server) handlerNotes(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method != http.MethodPost {
		fmt.Println("method OPTIONS in handler notes")
		return
	}

	data := map[string]string{
		"result": "access TRUE",
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

}

func (s *server) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method != http.MethodPost {
		fmt.Println("method OPTIONS in handler notes")
		return
	}

	type req struct {
		Data struct {
			RefreshToken string `json:"refreshToken"`
			Fingerprint  string `json:"fingerprint"`
		} `json:"data"`
	}
	data := &req{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	refSession, err := s.store.AuthRepository.UpdateTokens(data.Data.RefreshToken, data.Data.Fingerprint)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	// cookie := &http.Cookie{
	// 	HttpOnly: true,
	// 	Name:     "refreshToken",
	// 	Value:    refSession.RefreshToken,
	// }
	// http.SetCookie(w, cookie)

	m := map[string]string{
		"refreshToken": refSession.RefreshToken,
		"accessToken":  refSession.AccessToken,
		"exp":          refSession.Exp,
	}

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		//logger.NewLog("server", "handlerCreateUser", err, m, 6, "http response: json encode error")
		return
	}

}

func (s *server) error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	s.respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func (s *server) middlewareNoCors() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			if r.Method == "OPTIONS" {
				fmt.Println("OPTIONS")
				return
			}
			f(w, r)
		}
	}
}

func (s *server) middlewareAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				return
			}

			type data struct {
				Data struct {
					AccessToken string `json:"accessToken"`
				} `json:"data"`
			}
			d := &data{}

			err := json.NewDecoder(r.Body).Decode(d)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			err = s.store.AuthRepository.VerifyAccessToken(d.Data.AccessToken)
			if err != nil {
				s.error(w, r, http.StatusUnauthorized, err)
				return
			}

			f(w, r)
		}
	}
}
