package service

import (
	"noteapp/internal/model"
	"noteapp/pkg/logger"
	"strconv"
)

type NotesRepository interface {
	// GROUPS
	AddGroup(email string, nameGroup string, pid int) error
	DelGroup(id int, email string) error
	UpdateGroup(id int, email string, newNameGroup string) error
	// NOTES
	AddNote(email string, title string, text string, group_id int) error
	DelNote(id int, email string) error
	UpdateNote(id int, email string, title string, text string, group_id int) error
	GetNotesList(email string) (model.NoteList, error)
	GetNote(id int, email string) (model.Note, error)
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

func (s *NotesService) AddGroup(email string, nameGroup string, pid int) error {

	err := s.repository.AddGroup(email, nameGroup, pid)
	if err != nil {
		logger.NewLog("service - AddGroup()", 2, err, "Filed to add group in repository", nil)
	}
	return err
}

func (s *NotesService) DelGroup(id int, email string) error {
	err := s.repository.DelGroup(id, email)
	if err != nil {
		logger.NewLog("service - DelGroup()", 2, err, "Filed to del group in repository", nil)
	}
	return err
}

func (s *NotesService) UpdateGroup(id int, email string, newNameGroup string) error {
	err := s.repository.UpdateGroup(id, email, newNameGroup)
	if err != nil {
		logger.NewLog("service - UpdateGroup()", 2, err, "Filed to update group in repository", nil)
	}
	return err
}

// NOTES

func (s *NotesService) AddNote(email string, title string, text string, group_id int) error {
	err := s.repository.AddNote(email, title, text, group_id)
	if err != nil {
		gID := strconv.Itoa(group_id)
		m := map[string]string{
			"email":    email,
			"title":    title,
			"text":     text,
			"group_id": gID,
		}
		logger.NewLog("service - AddNote()", 2, err, "Filed to add note in repository", m)
	}
	return err
}

func (s *NotesService) DelNote(id int, email string) error {
	err := s.repository.DelNote(id, email)
	if err != nil {
		nID := strconv.Itoa(id)
		m := map[string]string{
			"note_id": nID,
		}
		logger.NewLog("service - DelNote()", 2, err, "Filed to del note in repository", m)
	}
	return err
}

func (s *NotesService) UpdateNote(id int, email string, title string, text string, group_id int) error {
	err := s.repository.UpdateNote(id, email, title, text, group_id)
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

func (s *NotesService) GetNotesList(email string) (model.NoteList, error) {
	list, err := s.repository.GetNotesList(email)
	if err != nil {
		logger.NewLog("service - GetNotesList()", 2, err, "Filed get notes list in repository, email = "+email, nil)
	}
	return list, err
}

func (s *NotesService) GetNote(id int, email string) (model.Note, error) {
	note, err := s.repository.GetNote(id, email)
	if err != nil {
		nID := strconv.Itoa(id)
		m := map[string]string{
			"note_id": nID,
		}
		logger.NewLog("service - UpdateNote()", 2, err, "Filed to get note in repository", m)
	}
	return note, err
}
