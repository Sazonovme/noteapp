package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"noteapp/internal/service"
	"noteapp/pkg/logger"
	"strings"
	"time"
)

var (
	errHeaderAuthorizationNotExist = errors.New("expected header authorization")
	errInvalidHeaderAuthorization  = errors.New("invalid authorization header")
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
			w.Header().Set("Access-Control-Allow-Methods", "*")
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
			accessTokenArr, ok := r.Header["Authorization"]

			if !ok {
				logger.NewLog("api - middlewareAuth()", 2, errHeaderAuthorizationNotExist, "header authorization not exist", nil)
				apiError(w, r, http.StatusUnauthorized, errHeaderAuthorizationNotExist)
				return

			}
			if len(accessTokenArr) != 1 {
				logger.NewLog("api - middlewareAuth()",
					2, errInvalidHeaderAuthorization,
					"header must have only 1 value",
					map[string]interface{}{
						"len":  len(accessTokenArr),
						"data": accessTokenArr,
					})
				apiError(w, r, http.StatusUnauthorized, errInvalidHeaderAuthorization)
				return
			}

			accessTokenArr = strings.Split(accessTokenArr[0], " ")
			if len(accessTokenArr) != 2 || accessTokenArr[0] != "Bearer" {
				logger.NewLog("api - middlewareAuth()",
					2, errInvalidHeaderAuthorization,
					"signature header must be: Bearer token",
					map[string]interface{}{
						"len":  len(accessTokenArr),
						"data": accessTokenArr,
					})
				apiError(w, r, http.StatusUnauthorized, errInvalidHeaderAuthorization)
				return
			}

			login, err := service.VerifyAccessToken(accessTokenArr[1])
			if err != nil {
				if err == service.ErrTokenInvalid {
					apiError(w, r, http.StatusUnauthorized, err)
					return
				}
				apiError(w, r, http.StatusInternalServerError, nil)
				return
			}

			m := map[string]string{}
			if r.Method != http.MethodGet {
				err := json.NewDecoder(r.Body).Decode(&m)
				if err != nil {
					logger.NewLog("api - middlewareAuth()", 2, err, "Field to decode r.Body", nil)
					apiError(w, r, http.StatusInternalServerError, nil)
					return
				}
			}

			m["login"] = login

			ctx := context.WithValue(r.Context(), ctxKey{}, m)

			f(w, r.WithContext(ctx))
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
