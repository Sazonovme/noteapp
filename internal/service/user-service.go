package service

import (
	"errors"
	"noteapp/internal/model"
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
		return errUserExist
	}

	err = u.ValidateBeforeCreate()
	if err != nil {
		return err
	}

	if _, err = u.EncryptPassword(); err != nil {
		return err
	}

	err = s.repository.CreateUser(u)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) FindByLogin(login string) (*model.User, error) {
	return s.repository.FindByLogin(login)
}
