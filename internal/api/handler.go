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
	errWrongPassword         = errors.New("wrong email or password")
	errRequiredFieldsMissing = errors.New("required fields are missing or not filled in")
)

type ctxKey struct{}

type UserService interface {
	CreateUser(*model.User) error
	FindByLogin(string) (*model.User, error)
}

type NotesService interface {
	// GROUPS
	AddGroup(email string, nameGroup string, pid int) error
	DelGroup(id int, email string) error
	UpdateGroup(id int, email string, newNameGroup string, pid int) error
	// NOTES
	AddNote(email string, title string, group_id int) error
	DelNote(id int, email string) error
	UpdateNote(data map[string]string) error
	GetNotesList(email string) (model.NoteList, error)
	GetNote(id int, email string) (model.Note, error)
}

type Handler struct {
	UserService  UserService
	NotesService NotesService
}

func NewHandler(userService UserService, notesService NotesService) *Handler {
	return &Handler{
		UserService:  userService,
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

	//LOG OUT

	router.HandleFunc("/logout", chainMiddleware(
		h.logOut,
		middlewareAuth(),
		middlewareNoCors(),
		middlewareLogIn()),
	)

	// GROUPS

	// router.HandleFunc("/addGroup", chainMiddleware(
	// 	h.addGroup,
	// 	middlewareAuth(),
	// 	middlewareNoCors(),
	// 	middlewareLogIn()),
	// )

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

	if d.User.Email == "" || d.User.Password == "" || d.User.Fingerprint == "" {
		logger.NewLog("api - createUser()", 2, err, "email or password or fingerprint not found in r.Body", err)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = h.UserService.CreateUser(&d.User)
	if err == service.ErrUserExist ||
		err == model.ErrValidationPassword ||
		err == model.ErrValidationEmail {

		apiError(w, r, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	// AT ONCE AUTH
	refSession, err := service.MakeRefreshSession(d.User.Email, d.User.Fingerprint)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	respData := &service.RequestTokenData{
		AccessToken:  refSession.AccessToken,
		RefreshToken: refSession.RefreshToken,
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(respData)
	if err != nil {
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	logger.NewLog("api - CreateUser()", 5, nil,
		"OUT - User created and authorized "+time.Now().Format("02.01 15:04:05"), nil)
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

	if d.User.Email == "" || d.User.Password == "" || d.User.Fingerprint == "" {
		logger.NewLog("api - authUser()", 2, err, "email or password or fingerprint not found in r.Body", err)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	u, err := h.UserService.FindByLogin(d.User.Email)
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

	refSession, err := service.MakeRefreshSession(u.Email, d.User.Fingerprint)
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

	if reqData.Data.Fingerprint == "" || reqData.Data.RefreshToken == "" {
		logger.NewLog("api - refreshToken()", 2, err, "refreshtoken or fingerprint not found in r.Body", err)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	refSession, err := service.UpdateTokens(reqData.Data.RefreshToken, reqData.Data.Fingerprint)
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

// LOG OUT

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
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

	refSession, err := service.UpdateTokens(reqData.Data.RefreshToken, reqData.Data.Fingerprint)
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

	email, ok1 := data["email"]
	name, ok2 := data["name"]
	if !(ok1 && ok2 && email != "" && name != "") {
		logger.NewLog("api - addNote()", 2, nil, "Required fields are missing in r.Context",
			"email = "+email+" name = "+name)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	pid := 0
	string_pid, ok3 := data["parentGroupId"]
	if ok3 && string_pid != "" {
		pid2, err := strconv.Atoi(string_pid)
		if err != nil {
			logger.NewLog("api - addGroup()", 2, err, "Filed to convert string to int", "string = "+string_pid)
			apiError(w, r, http.StatusInternalServerError, nil)
			return
		}
		pid = pid2
	}

	err := h.NotesService.AddGroup(email, name, pid)
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

	id_string := r.URL.Query().Get("id")
	email, ok1 := data["email"]
	if !(ok1 && email != "" && id_string != "") {
		logger.NewLog("api - delGroup()", 2, nil, "Required fields are missing in r.Context",
			"email = "+email+" id = "+id_string)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		logger.NewLog("api - delGroup()", 2, err, "Filed to convert string to int", "string = "+id_string)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	err = h.NotesService.DelGroup(id, email)
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
	email, ok2 := data["email"]
	if !(ok1 && ok2 && email != "" && id_string != "") {
		logger.NewLog("api - updateGroup()", 2, nil, "Required fields are missing in r.Context",
			"id = "+id_string+" email = "+email)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		logger.NewLog("api - updateGroup()", 2, err, "Filed to convert string to int", "string = "+id_string)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	name := data["name"]

	pid := -1
	string_pid, ok4 := data["parentGroupId"]
	if ok4 && string_pid != "" {
		pid2, err := strconv.Atoi(string_pid)
		if err != nil {
			logger.NewLog("api - updateGroup()", 2, err, "Filed to convert string to int", "string = "+string_pid)
			apiError(w, r, http.StatusInternalServerError, nil)
			return
		}
		pid = pid2
	}

	err = h.NotesService.UpdateGroup(id, email, name, pid)
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

	email, ok1 := data["email"]
	title, ok2 := data["title"]
	if !(ok1 && ok2 && email != "" && title != "") {
		logger.NewLog("api - addNote()", 2, nil, "Required fields are missing in r.Context", nil)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	var err error
	group_id := -1
	group_id_string, ok3 := data["group_id"]
	if ok3 && group_id_string != "" {
		group_id, err = strconv.Atoi(group_id_string)
		if err != nil {
			logger.NewLog("api - addNote()", 2, err, "Filed to convert string to int", "string = "+group_id_string)
			apiError(w, r, http.StatusInternalServerError, nil)
			return
		}
	}

	err = h.NotesService.AddNote(email, title, group_id)
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

	string_id := r.URL.Query().Get("id")
	email, ok1 := data["email"]
	if !(ok1 && string_id != "" && email != "") {
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

	err = h.NotesService.DelNote(id, email)
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
	email, ok2 := data["email"]
	if !(ok1 && ok2 && string_id != "" && email != "") {
		logger.NewLog("api - updateNote()", 2, nil, "Required fields are missing in r.Contex", nil)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	// id, err := strconv.Atoi(string_id)
	// if err != nil {
	// 	logger.NewLog("api - updateNote()", 2, err, "Filed to convert string to int", "string = "+string_id)
	// 	apiError(w, r, http.StatusInternalServerError, nil)
	// 	return
	// }

	// text := data["text"]

	// title := data["title"]

	// group_id_string := data["group_id"]

	// var group_id int
	// if group_id_string == "" {
	// 	group_id = -1
	// } else {
	// 	group_id, err = strconv.Atoi(group_id_string)
	// 	if err != nil {
	// 		logger.NewLog("api - updateNote()", 2, err, "Filed to convert string to int", "string = "+group_id_string)
	// 		apiError(w, r, http.StatusInternalServerError, nil)
	// 		return
	// 	}
	// }

	err := h.NotesService.UpdateNote(data)
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

	email, ok1 := data["email"]
	if !(ok1 && email != "") {
		logger.NewLog("api - getNotesList()", 2, nil, "Required fields are missing in r.Contex", "email")
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	list, err := h.NotesService.GetNotesList(email)
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

	email, ok := data["email"]
	string_id := r.URL.Query().Get("id")
	if !(ok && email != "" && string_id != "") {
		logger.NewLog("api - getNote()", 2, nil, "Required fields are missing in r.Contex", "email="+email+" id="+string_id)
		apiError(w, r, http.StatusBadRequest, errRequiredFieldsMissing)
		return
	}

	id, err := strconv.Atoi(string_id)
	if err != nil {
		logger.NewLog("api - getNote()", 2, err, "Filed to convert string to int", "string = "+string_id)
		apiError(w, r, http.StatusInternalServerError, nil)
		return
	}

	note, err := h.NotesService.GetNote(id, email)
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
