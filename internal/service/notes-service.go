package service

import (
	"noteapp/internal/repository"
	"noteapp/pkg/logger"
)

type NotesRepository interface {
	AddGroup(login string, nameGroup string) error
	DelGroup(login string, nameGroup string) error
	UpdateGroup(login string, id int, nameGroup string) error
	GetGroupList(login string) ([]repository.Group, error)
}

type NotesService struct {
	repository NotesRepository
}

func NewNotesService(repo NotesRepository) *NotesService {
	return &NotesService{
		repository: repo,
	}
}

// GROUPS

func (s *NotesService) AddGroup(login string, nameGroup string) error {
	err := s.repository.AddGroup(login, nameGroup)
	if err != nil {
		logger.NewLog("service - AddGroup()", 2, err, "Filed to add group in repository", nil)
	}
	return err
}
func (s *NotesService) DelGroup(login string, nameGroup string) error {
	err := s.repository.DelGroup(login, nameGroup)
	if err != nil {
		logger.NewLog("service - AddGroup()", 2, err, "Filed to del group in repository", nil)
	}
	return err
}
func (s *NotesService) UpdateGroup(login string, id int, newNameGroup string) error {
	err := s.repository.UpdateGroup(login, id, newNameGroup)
	if err != nil {
		logger.NewLog("service - AddGroup()", 2, err, "Filed to update group in repository", nil)
	}
	return err
}
func (s *NotesService) GetGroupList(login string) ([]repository.Group, error) {
	groupList, err := s.repository.GetGroupList(login)
	if err != nil {
		logger.NewLog("service - AddGroup()", 2, err, "Filed to get group list in repository", nil)
		return groupList, err
	}
	return groupList, nil
}

// NOTES
