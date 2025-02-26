package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword    = errors.New("invalid login or password")
	ErrValidationEmail    = errors.New("invalid email")
	ErrValidationPassword = errors.New("err password validation len - max = 50, min = 5, chars - ACSII")
)

type User struct {
	Id          int    `json:"-"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

func (u *User) ValidateBeforeCreate() error {

	// email
	ok1 := govalidator.IsEmail(u.Email)
	if !ok1 {
		return ErrValidationEmail
	}

	// password
	ok1 = govalidator.MaxStringLength(u.Password, "50")
	ok2 := govalidator.MinStringLength(u.Password, "5")
	ok3 := govalidator.IsASCII(u.Password)
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
