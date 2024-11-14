package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"noteapp/internal/model"
	"noteapp/internal/repository"
	"noteapp/internal/service"
	"noteapp/pkg/logger"
	"time"
)

var (
	errMethodNotAllowed = errors.New("method not allowed")
	errWrongPassword    = errors.New("wrong login or password")
)

type UserService interface {
	CreateUser(*model.User) error
	FindByLogin(string) (*model.User, error)
}

type AuthService interface {
	MakeRefreshSession(string, string) (*service.RequestTokenData, error)
	UpdateTokens(oldRefreshToken string, fingerprint string, login string) (*service.RequestTokenData, error)
}

type NotesService interface {
	AddGroup(login string, nameGroup string) error
	DelGroup(login string, nameGroup string) error
	UpdateGroup(login string, id int, nameGroup string) error
	GetGroupList(login string) ([]repository.Group, error)
}

type Handler struct {
	UserService  UserService
	AuthService  AuthService
	NotesService NotesService
}

func NewHandler(userService UserService, authService AuthService, notesService NotesService) *Handler {
	return &Handler{
		UserService:  userService,
		AuthService:  authService,
		NotesService: notesService,
	}
}

func (h *Handler) InitRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/sign-up", chainMiddleware(
		h.createUser,
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/sign-in", chainMiddleware(
		h.authUser,
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/refresh-token", chainMiddleware(
		h.RefreshToken,
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/addGroup", chainMiddleware(
		h.addGroup,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/delGroup", chainMiddleware(
		h.delGroup,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/updateGroup", chainMiddleware(
		h.updateGroup,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/getGroupList", chainMiddleware(
		h.getGroupList,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	return router
}

// USERS

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	d := struct {
		User model.User `json:"data"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		logger.NewLog("api - createUser()", 2, err, "Filed to decode r.Body", err)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = h.UserService.CreateUser(&d.User)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	logger.NewLog("api - CreateUser()", 5, nil,
		"OUT - User created "+time.Now().Format("02.01 15:04:05"), nil)
}

// AUTH

func (h *Handler) authUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	d := struct {
		User model.User `json:"data"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	u, err := h.UserService.FindByLogin(d.User.Login)
	if err != nil && err == repository.ErrUserNotFound {
		apiError(w, r, http.StatusBadRequest, repository.ErrUserNotFound)
		return
	} else if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = u.ComparePassword(d.User.Password)
	if err != nil {
		apiError(w, r, http.StatusUnauthorized, errWrongPassword)
		return
	}

	refSession, err := h.AuthService.MakeRefreshSession(u.Login, u.Fingerprint)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	// cookie := &http.Cookie{
	// 	HttpOnly: true,
	// 	Name:     "refreshToken",
	// 	Value:    refSession.RefreshToken,
	// }
	// http.SetCookie(w, cookie)

	respData := &service.RequestTokenData{
		AccessToken:  refSession.AccessToken,
		RefreshToken: refSession.RefreshToken,
		Exp:          refSession.Exp,
	}
	err = json.NewEncoder(w).Encode(respData)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}
	logger.NewLog("api - AuthUser", 5, nil,
		"OUT - User is authorized "+time.Now().Format("02.01 15:04:05"), respData)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	reqData := &struct {
		Data struct {
			RefreshToken string `json:"refreshToken"`
			Fingerprint  string `json:"fingerprint"`
			Login        string `json:"login"`
		} `json:"data"`
	}{}
	err := json.NewDecoder(r.Body).Decode(reqData)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	refSession, err := h.AuthService.UpdateTokens(reqData.Data.RefreshToken, reqData.Data.Fingerprint, reqData.Data.Login)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	// cookie := &http.Cookie{
	// 	HttpOnly: true,
	// 	Name:     "refreshToken",
	// 	Value:    refSession.RefreshToken,
	// }
	// http.SetCookie(w, cookie)

	err = json.NewEncoder(w).Encode(refSession)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - RefreshToken()", 5, nil,
		"OUT - New tokens generated "+time.Now().Format("02.01 15:04:05"), refSession)
}

// NOTES GROUPS

func (h *Handler) addGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	var group model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		logger.NewLog("api - addGroup()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.AddGroup(group.User_login, group.Name)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	logger.NewLog("api - addGroup()", 5, nil,
		"OUT - Group added "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) delGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	var group model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		logger.NewLog("api - delGroup()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.DelGroup(group.User_login, group.Name)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - delGroup()", 5, nil,
		"OUT - Group deleted "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) updateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	var group model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		logger.NewLog("api - updateGroup()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.UpdateGroup(group.User_login, group.Id, group.Name)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - updateGroup()", 5, nil,
		"OUT - Group updated "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) getGroupList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	var group model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		logger.NewLog("api - getGroupList()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	gList, err := h.NotesService.GetGroupList(group.User_login)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	if err := json.NewEncoder(w).Encode(gList); err != nil {
		logger.NewLog("api - getGroupList()", 2, err, "Filed to encode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - getGroupList()", 5, nil,
		"OUT - List geted "+time.Now().Format("02.01 15:04:05"), nil)
}
