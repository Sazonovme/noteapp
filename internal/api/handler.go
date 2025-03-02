package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"noteapp/internal/model"
	"noteapp/internal/repository"
	"noteapp/internal/service"
	"noteapp/pkg/logger"
	"strconv"
	"time"
)

var (
	errMethodNotAllowed      = errors.New("method not allowed")
	errWrongPassword         = errors.New("wrong login or password")
	errRequiredFieldsMissing = errors.New("required fields are missing or not filled in")
)

type ctxKey struct{}

type UserService interface {
	CreateUser(*model.User) error
	FindByLogin(string) (*model.User, error)
}

type AuthService interface {
	MakeRefreshSession(string, string) (*service.RequestTokenData, error)
	UpdateTokens(oldRefreshToken string, fingerprint string) (*service.RequestTokenData, error)
}

type NotesService interface {
	// GROUPS
	AddGroup(login string, nameGroup string) error
	DelGroup(id int, login string) error
	UpdateGroup(id int, login string, newNameGroup string) error
	GetGroupList(login string) (model.GroupList, error)
	// NOTES
	AddNote(login string, title string, text string, group_id int) error
	DelNote(id int, login string) error
	UpdateNote(id int, login string, title string, text string, group_id int) error
	GetNotesList(login string, group_id int) (model.NoteList, error)
	GetNote(id int, login string) (model.Note, error)
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

func (h *Handler) InitHandler() *http.ServeMux {
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
	if err == service.ErrUserExist ||
		err == model.ErrValidationPassword ||
		err == model.ErrValidationLogin {

		apiError(w, r, http.StatusBadRequest, err)
		return
	}

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

	respData := &service.RequestTokenData{
		AccessToken:  refSession.AccessToken,
		RefreshToken: refSession.RefreshToken,
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

	refSession, err := h.AuthService.UpdateTokens(reqData.Data.RefreshToken, reqData.Data.Fingerprint)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - addGroup()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok1 := data["login"]
	name, ok2 := data["name"]
	if !(ok1 && ok2 && login != "" && name != "") {
		logger.NewLog("api - addNote()", 2, nil, "Required fields are missing in r.Context",
			"login = "+login+"name = "+name)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	err := h.NotesService.AddGroup(login, name)
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - delGroup()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	id_string, ok1 := data["id"]
	login, ok2 := data["login"]
	if !(ok1 && ok2 && login != "" && id_string != "") {
		logger.NewLog("api - delGroup()", 2, nil, "Required fields are missing in r.Context",
			"login = "+login+"id = "+id_string)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		logger.NewLog("api - delGroup()", 2, err, "Filed to convert string to int", "string = "+id_string)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = h.NotesService.DelGroup(id, login)
	if err == repository.ErrInvalidData {
		apiError(w, r, http.StatusBadRequest, repository.ErrInvalidData)
		return
	}
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - updateGroup()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	id_string, ok1 := data["id"]
	login, ok2 := data["login"]
	name, ok3 := data["name"]
	if !(ok1 && ok2 && ok3 && login != "" && name != "" && id_string != "") {
		logger.NewLog("api - updateGroup()", 2, nil, "Required fields are missing in r.Context",
			"id = "+id_string+" login = "+login+" name = "+name)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		logger.NewLog("api - updateGroup()", 2, err, "Filed to convert string to int", "string = "+id_string)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = h.NotesService.UpdateGroup(id, login, name)
	if err == repository.ErrInvalidData {
		apiError(w, r, http.StatusBadRequest, repository.ErrInvalidData)
		return
	}
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - getGroupList()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := data["login"]
	if !(ok && login != "") {
		logger.NewLog("api - getGroupList()", 2, nil, "Field login not exist in r.Context()", "login = "+login)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - addNote()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok1 := data["login"]
	title, ok2 := data["title"]
	text, ok3 := data["text"]
	group_id_string, ok4 := data["group_id"]
	if !(ok1 && ok2 && ok3 && ok4 && login != "" && title != "") {
		logger.NewLog("api - addNote()", 2, nil, "Required fields are missing in r.Context", nil)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	var group_id int
	var err error
	if group_id_string == "" {
		group_id = 0
	} else {
		group_id, err = strconv.Atoi(group_id_string)
		if err != nil {
			logger.NewLog("api - addNote()", 2, err, "Filed to convert string to int", "string = "+group_id_string)
			apiError(w, r, http.StatusInternalServerError, nil)
			return
		}
	}

	err = h.NotesService.AddNote(login, title, text, group_id)
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - delNote()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	string_id, ok1 := data["id"]
	login, ok2 := data["login"]
	if !(ok1 && ok2 && string_id != "" && login != "") {
		logger.NewLog("api - delNote()", 2, nil, "Required fields are missing in r.Contex", nil)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(string_id)
	if err != nil {
		logger.NewLog("api - delNote()", 2, err, "Filed to convert string to int", "string = "+string_id)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = h.NotesService.DelNote(id, login)
	if err == repository.ErrInvalidData {
		apiError(w, r, http.StatusBadRequest, repository.ErrInvalidData)
		return
	}
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - updateNote()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	string_id, ok1 := data["id"]
	login, ok2 := data["login"]
	title, ok3 := data["title"]
	text, ok4 := data["text"]
	group_id_string, ok5 := data["group_id"]
	if !(ok1 && ok2 && ok3 && ok4 && ok5 && string_id != "" && login != "" && title != "") {
		logger.NewLog("api - updateNote()", 2, nil, "Required fields are missing in r.Contex", nil)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(string_id)
	if err != nil {
		logger.NewLog("api - updateNote()", 2, err, "Filed to convert string to int", "string = "+string_id)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	var group_id int
	if group_id_string == "" {
		group_id = 0
	} else {
		group_id, err = strconv.Atoi(group_id_string)
		if err != nil {
			logger.NewLog("api - updateNote()", 2, err, "Filed to convert string to int", "string = "+group_id_string)
			apiError(w, r, http.StatusInternalServerError, nil)
			return
		}
	}

	err = h.NotesService.UpdateNote(id, login, title, text, group_id)
	if err == repository.ErrInvalidData {
		apiError(w, r, http.StatusBadRequest, repository.ErrInvalidData)
		return
	}
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - getNotesList()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok1 := data["login"]
	if !(ok1 && login != "") {
		logger.NewLog("api - getNotesList()", 2, nil, "Required fields are missing in r.Contex", "login")
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	group_id_string := r.URL.Query().Get("group_id")
	var group_id int
	var err error
	if group_id_string == "" {
		group_id = 0
	} else {
		group_id, err = strconv.Atoi(group_id_string)
		if err != nil {
			logger.NewLog("api - getNotesList()", 2, err, "Filed to convert string to int", "string = "+group_id_string)
			apiError(w, r, http.StatusInternalServerError, nil)
			return
		}
	}

	list, err := h.NotesService.GetNotesList(login, group_id)
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

	data, ok := r.Context().Value(ctxKey{}).(map[string]string)
	if !ok {
		logger.NewLog("api - getNote()", 2, nil, "Filed to recive contex data", nil)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	login, ok := data["login"]
	string_id := r.URL.Query().Get("id")
	if !(ok && login != "" && string_id != "") {
		logger.NewLog("api - getNote()", 2, nil, "Required fields are missing in r.Contex", "login="+login+" id="+string_id)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(string_id)
	if err != nil {
		logger.NewLog("api - getNote()", 2, err, "Filed to convert string to int", "string = "+string_id)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	note, err := h.NotesService.GetNote(id, login)
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

func (h *Handler) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем query параметры
	noteIDParam := r.URL.Query().Get("note_id")
	login := r.URL.Query().Get("login")

	// Проверяем, что параметры переданы
	if noteIDParam == "" || login == "" {
		http.Error(w, "Missing note_id or login parameter", http.StatusBadRequest)
		return
	}

	// Преобразуем note_id в int
	noteID, err := strconv.Atoi(noteIDParam)
	if err != nil {
		http.Error(w, "Invalid note_id parameter", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для удаления
	err = h.NotesService.DelNote(noteID, login)
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
}
