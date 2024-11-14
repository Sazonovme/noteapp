package repository

import (
	"database/sql"
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

func (r *NotesRepository) AddGroup(login string, nameGroup string) error {
	_, err := r.db.Exec("INSERT INTO groups(user_login, name) VALUES ($1, $2)", login, nameGroup)
	return err
}
func (r *NotesRepository) DelGroup(login string, nameGroup string) error {
	_, err := r.db.Exec("DELETE FROM groups WHERE nameGroup = $1", nameGroup)
	return err
}
func (r *NotesRepository) UpdateGroup(login string, id int, newNameGroup string) error {
	_, err := r.db.Exec("UPDATE groups SET nameGroup = $1 WHERE id = $2 AND user_login = $3", newNameGroup, id, login)
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
