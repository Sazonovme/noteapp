package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/model"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ALL ERRORS
// errMethodNotAllowed      = errors.New("method not allowed")
// 	errWrongPassword         = errors.New("wrong login or password")
// 	errRequiredFieldsMissing = errors.New("required fields are missing or not filled in")
// ErrUserExist = errors.New             ("user with this login already exist")
// ErrInvalidPassword    = errors.New	("invalid login or password") // invalid password
// ErrValidationLogin    = errors.New	("err login validation len - max = 50, min = 5, chars - only letters and nums")
// ErrValidationPassword = errors.New	("err password validation len - max = 50, min = 5, chars - ACSII")

func TestCreateUser(t *testing.T) {

	type want struct {
		code     int
		response map[string]string
	}
	type payload struct {
		Login    string
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
				Login:    "validLogin",
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
				Login:    "invl",
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
				Login:    "invl12",
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
				Login:    "invl12",
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
				Login:    "validLogin",
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
					Login:       tcase.payload.Login,
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
