package repository

import (
	"database/sql"
	"noteapp/internal/model"
)

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NotesRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *NotesRepository {
	return &NotesRepository{
		db: db,
	}
}

// GROUPS

func (r *NotesRepository) AddGroup(login string, nameGroup string) error {
	_, err := r.db.Exec("INSERT INTO groups(user_login, name) VALUES ($1, $2)", login, nameGroup)
	return err
}

func (r *NotesRepository) DelGroup(id int, login string) error {
	_, err := r.db.Exec("DELETE FROM groups WHERE id = $1 AND user_login = $2", id)
	return err
}

func (r *NotesRepository) UpdateGroup(id int, login string, newNameGroup string) error {
	_, err := r.db.Exec("UPDATE groups SET name = $1 WHERE id = $2 AND user_login", newNameGroup, id, login)
	return err
}

func (r *NotesRepository) GetGroupList(login string) ([]Group, error) {
	var list []Group

	res, err := r.db.Query("SELECT id, name FROM groups WHERE login = $1", login)
	if err != nil {
		if err == sql.ErrNoRows {
			return list, nil
		}
		return list, err
	}
	defer res.Close()

	for res.Next() {
		var listElem Group
		if err := res.Scan(&listElem.Id, &listElem.Name); err != nil {
			return list, err
		}
		list = append(list, listElem)
	}

	if err = res.Err(); err != nil {
		return list, err
	}
	return list, nil
}

// NOTES

func (r *NotesRepository) AddNote(login string, title string, text string, group_id int) error {
	_, err := r.db.Exec("INSERT INTO notes(user_login, title, text, group_id) VALUES ($1, $2, $3, $4)", login, title, text, group_id)
	return err
}

func (r *NotesRepository) DelNote(id int, login string) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id = $1 AND login = $2", id, login)
	return err
}

func (r *NotesRepository) UpdateNote(id int, login string, title string, text string, group_id int) error {
	_, err := r.db.Exec("UPDATE notes SET title = $1, text = $2, group_id = $3 WHERE id = $4 AND user_login = $5", title, text, id, group_id, login)
	return err
}

func (r *NotesRepository) GetNotesList(login string, group_id int) ([]model.Note, error) {
	var res *sql.Rows
	var err error
	var list []model.Note

	if group_id != 0 {
		res, err = r.db.Query(
			"SELECT id, title, group_id FROM notes WHERE user_login = $1 AND group_id = $2",
			login,
			group_id,
		)
	} else {
		res, err = r.db.Query(
			"SELECT id, title, group_id FROM notes WHERE user_login = $1",
			login,
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return list, nil
		}
		return list, err
	}
	defer res.Close()

	for res.Next() {
		var listElem model.Note
		if err := res.Scan(&listElem.Id, &listElem.Title, &listElem.Group_id); err != nil {
			return list, err
		}
		list = append(list, listElem)
	}

	if err = res.Err(); err != nil {
		return list, err
	}
	return list, nil
}

func (r *NotesRepository) GetNote(id int, login string) (model.Note, error) {

	var note model.Note

	err := r.db.QueryRow(
		"SELECT id, title, text, group_id FROM notes WHERE id = $1 AND user_login = $2",
		id,
		login,
	).Scan(&note)
	if err != nil {
		return note, err
	}

	return note, nil
}
