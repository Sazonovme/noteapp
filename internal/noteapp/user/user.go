package user

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New() *User {
	return &User{}
}

func (u *User) BeforeCreate() error {
	// 1. Проверить существует ли такой
	// 2. Назначить ID
	//...
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
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
