package service

import (
	"errors"
	"noteapp/internal/model"
	"noteapp/pkg/logger"
)

var (
	errUserExist = errors.New("user with this login already exist")
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
	_, err := s.repository.FindByLogin(u.Login)
	if err == nil {
		logger.NewLog("service - CreateUser()", 5, err, "user exist", u.Login)
		return errUserExist
	}

	err = u.ValidateBeforeCreate()
	if err != nil {
		logger.NewLog("service - CreateUser()", 5, err, "Filed to validate before create", u.Login)
		return err
	}

	if _, err = u.EncryptPassword(); err != nil {
		logger.NewLog("service - CreateUser()", 2, err, "Filed to encrypt password", u.Login)
		return err
	}

	err = s.repository.CreateUser(u)
	if err != nil {
		logger.NewLog("service - CreateUser()", 2, err, "Filed to create user in repository", u.Login)
		return err
	}

	return nil
}

func (s *UserService) FindByLogin(login string) (*model.User, error) {
	return s.repository.FindByLogin(login)
}
