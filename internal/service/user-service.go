package service

import (
	"errors"
	"noteapp/internal/model"
	"noteapp/pkg/logger"
)

var (
	ErrUserExist = errors.New("user with this email already exist")
)

type UserRepository interface {
	CreateUser(*model.User) error
	FindByLogin(string) (*model.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (s *UserService) CreateUser(u *model.User) error {
	_, err := s.repository.FindByLogin(u.Email)
	if err == nil {
		logger.NewLog("service - CreateUser()", 5, err, "user exist", u.Email)
		return ErrUserExist
	}

	err = u.ValidateBeforeCreate()
	if err != nil {
		logger.NewLog("service - CreateUser()", 5, err, "Filed to validate before create", u.Email)
		return err
	}

	if _, err = u.EncryptPassword(); err != nil {
		logger.NewLog("service - CreateUser()", 2, err, "Filed to encrypt password", u.Email)
		return err
	}

	err = s.repository.CreateUser(u)
	if err != nil {
		logger.NewLog("service - CreateUser()", 2, err, "Filed to create user in repository", u.Email)
		return err
	}

	return nil
}

func (s *UserService) FindByLogin(email string) (*model.User, error) {
	return s.repository.FindByLogin(email)
}
