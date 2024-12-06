package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword    = errors.New("invalid login or password")
	ErrValidationLogin    = errors.New("err login validation len - max = 50, min = 5, chars - only letters and nums")
	ErrValidationPassword = errors.New("err password validation len - max = 50, min = 5, chars - ACSII")
)

type User struct {
	Id          int    `json:"-"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

func (u *User) ValidateBeforeCreate() error {

	// login
	ok1 := govalidator.MaxStringLength(u.Login, "50")
	ok2 := govalidator.MinStringLength(u.Login, "5")
	ok3 := govalidator.IsAlphanumeric(u.Login)
	if !(ok1 && ok2 && ok3) {
		return ErrValidationLogin
	}

	// password
	ok1 = govalidator.MaxStringLength(u.Password, "50")
	ok2 = govalidator.MinStringLength(u.Password, "5")
	ok3 = govalidator.IsASCII(u.Password)
	if !(ok1 && ok2 && ok3) {
		return ErrValidationPassword
	}

	return nil
}

func (u *User) EncryptPassword() (*User, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(data)
	return u, nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}
