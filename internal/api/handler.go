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

type UserService interface {
	CreateUser(*model.User) error
	FindByLogin(string) (*model.User, error)
}

type AuthService interface {
	MakeRefreshSession(string, string) (*service.RequestTokenData, error)
	UpdateTokens(oldRefreshToken string, fingerprint string, login string) (*service.RequestTokenData, error)
}

type NotesService interface {
	// GROUPS
	AddGroup(login string, nameGroup string) error
	DelGroup(id int, login string) error
	UpdateGroup(id int, login string, newNameGroup string) error
	GetGroupList(login string) ([]model.Group, error)
	// NOTES
	AddNote(login string, title string, text string, group_id int) error
	DelNote(id int, login string) error
	UpdateNote(id int, login string, title string, text string, group_id int) error
	GetNotesList(login string, group_id int) ([]model.Note, error)
	GetNote(id int, login string) (model.Note, error)
}

type Handler struct {
	UserService  UserService
	AuthService  AuthService
	NotesService NotesService
}

type ctxKey string

var (
	errMethodNotAllowed = errors.New("method not allowed")
	errWrongPassword    = errors.New("wrong login or password")
)

var ctxUserLogin ctxKey = "login"

func NewHandler(userService UserService, authService AuthService, notesService NotesService) *Handler {
	return &Handler{
		UserService:  userService,
		AuthService:  authService,
		NotesService: notesService,
	}
}

func (h *Handler) InitRouter() *http.ServeMux {
	router := http.NewServeMux()

	// AUTH

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
		h.refreshToken,
		middlewareNoCors(),
		middlewareLogIn()),
	)

	// GROUPS

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

	// NOTES

	router.HandleFunc("/addNote", chainMiddleware(
		h.addNote,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/delNote", chainMiddleware(
		h.delNote,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/updateNote", chainMiddleware(
		h.updateNote,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/getNotesList", chainMiddleware(
		h.getNotesList,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	router.HandleFunc("/getNote", chainMiddleware(
		h.getNote,
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

func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	reqData := &struct {
		Data struct {
			RefreshToken string `json:"refreshToken"`
			Fingerprint  string `json:"fingerprint"`
		} `json:"data"`
	}{}
	err := json.NewDecoder(r.Body).Decode(reqData)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - refreshToken()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	refSession, err := h.AuthService.UpdateTokens(reqData.Data.RefreshToken, reqData.Data.Fingerprint, login)
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

	data := struct {
		Group model.Group `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - addGroup()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - addGroup()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.AddGroup(login, data.Group.Name)
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

	data := struct {
		Group model.Group `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - delGroup()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - delGroup()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.DelGroup(data.Group.Id, login)
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

	data := struct {
		Group model.Group `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - updateGroup()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - updateGroup()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.UpdateGroup(data.Group.Id, login, data.Group.Name)
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

	// data := struct {
	// 	Group model.Group `json:"data"`
	// }{}
	// if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
	// 	logger.NewLog("api - getGroupList()", 2, err, "Filed to decode r.Body", nil)
	// 	apiError(w, r, http.StatusInternalServerError, nil)
	// 	return
	// }

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - getGroupList()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	gList, err := h.NotesService.GetGroupList(login)
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

// NOTES

func (h *Handler) addNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	data := struct {
		Note model.Note `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - addNote()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - addNote()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.AddNote(login, data.Note.Title, data.Note.Text, data.Note.Group_id)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	logger.NewLog("api - addNote()", 5, nil,
		"OUT - Note added "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) delNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	data := struct {
		Data struct {
			Id int `json:"id"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - delNote()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - delNote()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.DelNote(data.Data.Id, login)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - delNote()", 5, nil,
		"OUT - Note deleted "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	data := struct {
		Note model.Note `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - updateNote()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - updateNote()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err := h.NotesService.UpdateNote(data.Note.Id, login, data.Note.Title, data.Note.Text, data.Note.Group_id)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - updateNote()", 5, nil,
		"OUT - Note updated "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) getNotesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	data := struct {
		Data struct {
			Group_id int `json:"group_id"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - getNotesList()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - getNotesList()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	list, err := h.NotesService.GetNotesList(login, data.Data.Group_id)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		logger.NewLog("api - getNotesList()", 2, err, "Filed to encode r.Body", list)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - getNotesList()", 5, nil,
		"OUT - List geted "+time.Now().Format("02.01 15:04:05"), nil)
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiError(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	data := struct {
		Data struct {
			Id int `json:"id"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.NewLog("api - getNote()", 2, err, "Filed to decode r.Body", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := r.Context().Value(ctxUserLogin).(string)
	if !ok {
		logger.NewLog("api - getNote()", 2, nil, "Field login not exist in r.Context()", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	note, err := h.NotesService.GetNote(data.Data.Id, login)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		logger.NewLog("api - getNote()", 2, err, "Filed to encode r.Body", note)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - getNote()", 5, nil,
		"OUT - Note geted "+time.Now().Format("02.01 15:04:05"), nil)
}
