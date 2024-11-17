package service

import (
	"noteapp/internal/model"
	"noteapp/internal/repository"
	"noteapp/pkg/logger"
	"strconv"
)

type NotesRepository interface {
	// GROUPS
	AddGroup(login string, nameGroup string) error
	DelGroup(id int, login string) error
	UpdateGroup(id int, login string, newNameGroup string) error
	GetGroupList(login string) ([]repository.Group, error)
	// NOTES
	AddNote(login string, title string, text string, group_id int) error
	DelNote(id int, login string) error
	UpdateNote(id int, login string, title string, text string, group_id int) error
	GetNotesList(login string, group_id int) ([]model.Note, error)
	GetNote(id int, login string) (model.Note, error)
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

func (s *NotesService) DelGroup(id int, login string) error {
	err := s.repository.DelGroup(id, login)
	if err != nil {
		logger.NewLog("service - DelGroup()", 2, err, "Filed to del group in repository", nil)
	}
	return err
}

func (s *NotesService) UpdateGroup(id int, login string, newNameGroup string) error {
	err := s.repository.UpdateGroup(id, login, newNameGroup)
	if err != nil {
		logger.NewLog("service - UpdateGroup()", 2, err, "Filed to update group in repository", nil)
	}
	return err
}

func (s *NotesService) GetGroupList(login string) ([]repository.Group, error) {
	groupList, err := s.repository.GetGroupList(login)
	if err != nil {
		logger.NewLog("service - GetGroupList()", 2, err, "Filed to get group list in repository", nil)
		return groupList, err
	}
	return groupList, nil
}

// NOTES

func (s *NotesService) AddNote(login string, title string, text string, group_id int) error {
	err := s.repository.AddNote(login, title, text, group_id)
	if err != nil {
		gID := strconv.Itoa(group_id)
		m := map[string]string{
			"login":    login,
			"title":    title,
			"text":     text,
			"group_id": gID,
		}
		logger.NewLog("service - AddNote()", 2, err, "Filed to add note in repository", m)
	}
	return err
}

func (s *NotesService) DelNote(id int, login string) error {
	err := s.repository.DelNote(id, login)
	if err != nil {
		nID := strconv.Itoa(id)
		m := map[string]string{
			"note_id": nID,
		}
		logger.NewLog("service - DelNote()", 2, err, "Filed to del note in repository", m)
	}
	return err
}

func (s *NotesService) UpdateNote(id int, login string, title string, text string, group_id int) error {
	err := s.repository.UpdateNote(id, login, title, text, group_id)
	if err != nil {
		gID := strconv.Itoa(group_id)
		nID := strconv.Itoa(id)
		m := map[string]string{
			"id":       nID,
			"title":    title,
			"text":     text,
			"group_id": gID,
		}
		logger.NewLog("service - UpdateNote()", 2, err, "Filed to update note in repository", m)
	}
	return err
}

func (s *NotesService) GetNotesList(login string, group_id int) ([]model.Note, error) {
	list, err := s.repository.GetNotesList(login, group_id)
	if err != nil {
		gID := strconv.Itoa(group_id)
		m := map[string]string{
			"login":    login,
			"group_id": gID,
		}
		logger.NewLog("service - GetNotesList()", 2, err, "Filed get notes list in repository", m)
	}
	return list, err
}

func (s *NotesService) GetNote(id int, login string) (model.Note, error) {
	note, err := s.repository.GetNote(id, login)
	if err != nil {
		nID := strconv.Itoa(id)
		m := map[string]string{
			"note_id": nID,
		}
		logger.NewLog("service - UpdateNote()", 2, err, "Filed to get note in repository", m)
	}
	return note, err
}
