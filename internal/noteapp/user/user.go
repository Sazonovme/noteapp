package user

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id          int    `json:"-"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

func New() *User {
	return &User{}
}

func (u *User) ValidateBeforeCreate() error {

	var errString string

	// login
	ok := govalidator.MaxStringLength(u.Login, "50")
	if !ok {
		errString += "invalid length login, max lenght - 50\n"
	}
	ok = govalidator.IsAlphanumeric(u.Login)
	if !ok {
		errString += "invalid string login, valid characters are letters and numbers\n"
	}

	// password
	ok = govalidator.MaxStringLength(u.Password, "50")
	if !ok {
		errString += "invalid length password, max lenght - 50\n"
	}
	ok = govalidator.IsASCII(u.Password)
	if !ok {
		errString += "invalid string password, valid characters ASCII\n"
	}

	if errString != "" {
		return errors.New(errString)
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
		return errors.New("invalid password")
	}
	return nil
}
