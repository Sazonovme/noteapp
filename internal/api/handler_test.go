package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/model"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ALL ERRORS
// errMethodNotAllowed      = errors.New("method not allowed")
// errWrongPassword         = errors.New("wrong login or password")
// errRequiredFieldsMissing = errors.New("required fields are missing or not filled in")
// ErrUserExist = errors.New             ("user with this login already exist")
// ErrInvalidPassword    = errors.New	("invalid login or password") // invalid password
// ErrValidationLogin    = errors.New	("err login validation len - max = 50, min = 5, chars - only letters and nums")
// ErrValidationPassword = errors.New	("err password validation len - max = 50, min = 5, chars - ACSII")

// HELPERS
func HelperCreateUser(t *testing.T, handler *http.ServeMux, email string, password string) {
	t.Helper()

	type request struct {
		model.User `json:"data"`
	}

	data := new(bytes.Buffer)
	err := json.NewEncoder(data).Encode(request{
		model.User{
			Email:    email,
			Password: password,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/sign-up", data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Accept", "application/json")

	handler.ServeHTTP(rec, req)

	if rec.Code != 201 {
		if rec.Body.Len() != 0 {
			result := map[string]string{}
			err = json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatal("Decode err: " + err.Error())
			}
			t.Fatal(result)
		} else {
			t.Fatal("Err: body is clear and code = " + strconv.Itoa(rec.Code))
		}
	}
}

func HelperAuthUser(t *testing.T, handler *http.ServeMux, email string, password string) {
	t.Helper()

	type request struct {
		model.User `json:"data"`
	}
	type response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	data := new(bytes.Buffer)
	err := json.NewEncoder(data).Encode(request{
		model.User{
			Email:    email,
			Password: password,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/sign-in", data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Accept", "application/json")

	handler.ServeHTTP(rec, req)

	if rec.Code != 201 {
		if rec.Body.Len() != 0 {
			result := map[string]string{}
			err = json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatal("Decode err: " + err.Error())
			}
			t.Fatal(result)
		} else {
			t.Fatal("Err: body is clear and code = " + strconv.Itoa(rec.Code))
		}
	}
}

// USERS
func TestCreateUser(t *testing.T) {

	type want struct {
		code     int
		response map[string]string
	}
	type payload struct {
		Email    string
		Password string
	}
	type request struct {
		model.User `json:"data"`
	}

	testCases := []struct {
		name    string
		method  string
		payload payload
		want    want
	}{
		// test 1 valid
		{
			name:   "valid case",
			method: "POST",
			payload: payload{
				Email:    "validLogin",
				Password: "validPassword",
			},
			want: want{
				code:     201,
				response: map[string]string{},
			},
		},
		// test 2 invalid login
		{
			name:   "invalid login",
			method: "POST",
			payload: payload{
				Email:    "invl",
				Password: "123qwe",
			},
			want: want{
				code: 400,
				response: map[string]string{
					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
				},
			},
		},
		// test 3 invalid password
		{
			name:   "invalid password",
			method: "POST",
			payload: payload{
				Email:    "invl12",
				Password: "123",
			},
			want: want{
				code: 400,
				response: map[string]string{
					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
				},
			},
		},
		// test 4 invalid method
		{
			name:   "invalid method",
			method: "GET",
			payload: payload{
				Email:    "invl12",
				Password: "123123",
			},
			want: want{
				code: 405,
				response: map[string]string{
					"error": "method not allowed",
				},
			},
		},
		// test 5 exist user
		{
			name:   "user exist",
			method: "POST",
			payload: payload{
				Email:    "validLogin",
				Password: "123123",
			},
			want: want{
				code: 400,
				response: map[string]string{
					"error": "user with this login already exist",
				},
			},
		},
	}

	handler, create, teardown, db := NewTestHandler(t)

	create(db, t, "test_users")
	defer teardown(db, t, "test_users")

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			data := new(bytes.Buffer)
			err := json.NewEncoder(data).Encode(request{
				model.User{
					Email:       tcase.payload.Email,
					Password:    tcase.payload.Password,
					Fingerprint: "test-fingerprint",
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Access-Control-Allow-Origin", "*")
			req.Header.Set("Accept", "application/json")

			handler.ServeHTTP(rec, req)

			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
				assert.Equal(t, tcase.want.code, rec.Code)
				return
			}

			result := map[string]string{}
			err = json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatal("Decode err: " + err.Error())
			}

			assert.Equal(t, tcase.want.code, rec.Code)
			assert.Equal(t, tcase.want.response, result)
		})
	}
}

// AUTH
func TestAuthUser(t *testing.T) {

	type want struct {
		code          int
		response      map[string]string
		checkOnlyKeys bool
	}
	type payload struct {
		Email       string
		Password    string
		Fingerprint string
	}
	type request struct {
		model.User `json:"data"`
	}

	testCases := []struct {
		name    string
		method  string
		payload payload
		want    want
	}{
		// test 1 valid
		{
			name:   "valid case",
			method: "POST",
			payload: payload{
				Email:       "existUser",
				Password:    "password",
				Fingerprint: "fingerprint",
			},
			want: want{
				code: 200,
				response: map[string]string{
					"accessToken":  "",
					"refreshToken": "",
				},
				checkOnlyKeys: true,
			},
		},
		// test 2 valid another fingerprint
		{
			name:   "valid case anther fingerprint",
			method: "POST",
			payload: payload{
				Email:       "existUser",
				Password:    "password",
				Fingerprint: "fingerprint2",
			},
			want: want{
				code: 200,
				response: map[string]string{
					"accessToken":  "",
					"refreshToken": "",
				},
				checkOnlyKeys: true,
			},
		},
		// test 3 invalid login
		{
			name:   "invalid login",
			method: "POST",
			payload: payload{
				Email:       "invl",
				Password:    "password",
				Fingerprint: "fingerprint",
			},
			want: want{
				code: 400,
				response: map[string]string{
					"error": "user not found",
				},
				checkOnlyKeys: false,
			},
		},
		// test 4 invalid password
		{
			name:   "invalid password",
			method: "POST",
			payload: payload{
				Email:       "existUser",
				Password:    "invalidpass",
				Fingerprint: "fingerprint",
			},
			want: want{
				code: 401,
				response: map[string]string{
					"error": "wrong login or password",
				},
				checkOnlyKeys: false,
			},
		},
		// test 5 invalid method
		{
			name:   "invalid method",
			method: "GET",
			payload: payload{
				Email:       "existUser",
				Password:    "password",
				Fingerprint: "fingerprint",
			},
			want: want{
				code: 405,
				response: map[string]string{
					"error": "method not allowed",
				},
				checkOnlyKeys: false,
			},
		},
	}

	handler, create, teardown, db := NewTestHandler(t)

	create(db, t, "test_refreshsessions")
	defer teardown(db, t, "test_refreshsessions")

	HelperCreateUser(t, handler, "existUser", "password")

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			data := new(bytes.Buffer)
			err := json.NewEncoder(data).Encode(request{
				model.User{
					Email:       tcase.payload.Email,
					Password:    tcase.payload.Password,
					Fingerprint: tcase.payload.Fingerprint,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tcase.method, "/sign-in", data)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Access-Control-Allow-Origin", "*")
			req.Header.Set("Accept", "application/json")

			handler.ServeHTTP(rec, req)

			if rec.Body.Len() == 0 {
				assert.Equal(t, tcase.want.code, rec.Code)
				return
			}

			result := map[string]string{}
			err = json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatal("Decode err: " + err.Error())
			}

			//check required fields (keys)
			for wantKey := range tcase.want.response {
				_, ok := result[wantKey]
				if !ok {
					t.Fatal("Not found want key in result: " + wantKey)
				}
			}

			if !tcase.want.checkOnlyKeys {
				assert.Equal(t, tcase.want.response, result)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {

	type Data struct {
		RefreshToken string `json:"refreshToken"`
		Fingerprint  string `json:"fingerprint"`
	}
	type request struct {
		Data Data `json:"data"`
	}
	type response struct {
		AccessToken   string `json:"accessToken"`
		RefreshToken  string `json:"refreshToken"`
		checkOnlyKeys bool
	}

	testCases := []struct {
		name    string
		method  string
		payload request
		want    response
	}{
		// test 1 valid
		{
			name:   "valid case",
			method: "POST",
			payload: request{
				Data: Data{
					RefreshToken: "",
					Fingerprint:  "test_fingerprint",
				},
			},
			want: response{
				AccessToken:   "",
				RefreshToken:  "",
				checkOnlyKeys: true,
			},
		},
	}

	handler, create, teardown, db := NewTestHandler(t)

	create(db, t, "test_refreshsessions")
	defer teardown(db, t, "test_refreshsessions")

	// CREATE USER AND FIRST TOKEN
	HelperCreateUser(t, handler, "existUser", "password")

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {

			// encode request
			data := new(bytes.Buffer)
			err := json.NewEncoder(data).Encode(request{
				Data{
					RefreshToken: "",
					Fingerprint:  "",
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			// 			rec := httptest.NewRecorder()
			// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
			// 			req.Header.Set("Content-Type", "application/json")
			// 			req.Header.Set("Access-Control-Allow-Origin", "*")
			// 			req.Header.Set("Accept", "application/json")

			// 			handler.ServeHTTP(rec, req)

			// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
			// 				assert.Equal(t, tcase.want.code, rec.Code)
			// 				return
			// 			}

			// 			result := map[string]string{}
			// 			err = json.NewDecoder(rec.Body).Decode(&result)
			// 			if err != nil {
			// 				t.Fatal("Decode err: " + err.Error())
			// 			}

			// 			assert.Equal(t, tcase.want.code, rec.Code)
			// 			assert.Equal(t, tcase.want.response, result)
		})
	}
}

// NOTES GROUPS

func TestAddGroup(t *testing.T) {

	type request struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	type response struct {
		code     int
		response map[string]string
	}

	testCases := []struct {
		name    string
		method  string
		payload request
		want    response
	}{
		// test 1 valid
		{
			name:   "valid case",
			method: "POST",
			payload: request{
				Email: "existUser",
				Name:  "testGroup",
			},
			want: response{
				code:     201,
				response: map[string]string{},
			},
		},
	}

	handler, create, teardown, db := NewTestHandler(t)

	create(db, t, "test_groups")
	HelperCreateUser(t, handler, "existUser", "secretPassword")
	defer teardown(db, t, "test_groups")

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			data := new(bytes.Buffer)
			err := json.NewEncoder(data).Encode(request{
				Email: tcase.payload.Email,
				Name:  tcase.payload.Name,
			})
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tcase.method, "/addGroup", data)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Access-Control-Allow-Origin", "*")
			req.Header.Set("Accept", "application/json")

			handler.ServeHTTP(rec, req)

			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
				assert.Equal(t, tcase.want.code, rec.Code)
				return
			}

			result := map[string]string{}
			err = json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatal("Decode err: " + err.Error())
			}

			assert.Equal(t, tcase.want.code, rec.Code)
			assert.Equal(t, tcase.want.response, result)
		})
	}
}

// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestDelGroup(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestUpdateGroup(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestGetGroupList(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// // NOTES

// func TestAddNote(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestDelNote(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestUpdateNote(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestGetNotesList(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }

// func TestGetNote(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response map[string]string
// 	}
// 	type payload struct {
// 		Login    string
// 		Password string
// 	}
// 	type request struct {
// 		model.User `json:"data"`
// 	}

// 	testCases := []struct {
// 		name    string
// 		method  string
// 		payload payload
// 		want    want
// 	}{
// 		// test 1 valid
// 		{
// 			name:   "valid case",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "validPassword",
// 			},
// 			want: want{
// 				code:     201,
// 				response: map[string]string{},
// 			},
// 		},
// 		// test 2 invalid login
// 		{
// 			name:   "invalid login",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl",
// 				Password: "123qwe",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err login validation len - max = 50, min = 5, chars - only letters and nums",
// 				},
// 			},
// 		},
// 		// test 3 invalid password
// 		{
// 			name:   "invalid password",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "err password validation len - max = 50, min = 5, chars - ACSII",
// 				},
// 			},
// 		},
// 		// test 4 invalid method
// 		{
// 			name:   "invalid method",
// 			method: "GET",
// 			payload: payload{
// 				Login:    "invl12",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 405,
// 				response: map[string]string{
// 					"error": "method not allowed",
// 				},
// 			},
// 		},
// 		// test 5 exist user
// 		{
// 			name:   "user exist",
// 			method: "POST",
// 			payload: payload{
// 				Login:    "validLogin",
// 				Password: "123123",
// 			},
// 			want: want{
// 				code: 400,
// 				response: map[string]string{
// 					"error": "user with this login already exist",
// 				},
// 			},
// 		},
// 	}

// 	handler, create, teardown, db := NewTestHandler(t)

// 	create(db, t, "test_users")
// 	defer teardown(db, t, "test_users")

// 	for _, tcase := range testCases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			data := new(bytes.Buffer)
// 			err := json.NewEncoder(data).Encode(request{
// 				model.User{
// 					Login:       tcase.payload.Login,
// 					Password:    tcase.payload.Password,
// 					Fingerprint: "test-fingerprint",
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			rec := httptest.NewRecorder()
// 			req, _ := http.NewRequest(tcase.method, "/sign-up", data)
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Access-Control-Allow-Origin", "*")
// 			req.Header.Set("Accept", "application/json")

// 			handler.ServeHTTP(rec, req)

// 			if rec.Body.Len() == 0 && reflect.DeepEqual(tcase.want.response, map[string]string{}) {
// 				assert.Equal(t, tcase.want.code, rec.Code)
// 				return
// 			}

// 			result := map[string]string{}
// 			err = json.NewDecoder(rec.Body).Decode(&result)
// 			if err != nil {
// 				t.Fatal("Decode err: " + err.Error())
// 			}

// 			assert.Equal(t, tcase.want.code, rec.Code)
// 			assert.Equal(t, tcase.want.response, result)
// 		})
// 	}
// }
