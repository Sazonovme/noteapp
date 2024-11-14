package api

import (
	"encoding/json"
	"net/http"
	"noteapp/internal/service"
	"noteapp/pkg/logger"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func chainMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func middlewareNoCors() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			if r.Method == "OPTIONS" {
				logger.NewLog("api - middlewareNoCors()", 5, nil, "OUT - OPTIONS "+time.Now().Format("02.01 15:04:05"), nil)
				return
			}
			f(w, r)
		}
	}
}

func middlewareAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			d := &struct {
				Data struct {
					AccessToken string `json:"accessToken"`
				} `json:"data"`
			}{}
			err := json.NewDecoder(r.Body).Decode(d)
			if err != nil {
				apiError(w, r, http.StatusInternalServerError, nil)
				return
			}

			_, err = service.VerifyAccessToken(d.Data.AccessToken)
			if err != nil {
				http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
				//apiError(w, r, http.StatusInternalServerError, nil)
				return
			}

			f(w, r)
		}
	}
}

func middlewareLogIn() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.NewLog("api - middlewareLogIn()", 5, nil, "IN "+time.Now().Format("02.01 15:04:05"), nil)
			f(w, r)
		}
	}
}
