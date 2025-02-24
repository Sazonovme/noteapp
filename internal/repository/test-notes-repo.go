package repository

import (
	"database/sql"
	"noteapp/internal/model"
)

type TestNotesRepository struct {
	db *sql.DB
}

func NewTestNotesRepository(db *sql.DB) *TestNotesRepository {
	return &TestNotesRepository{
		db: db,
	}
}

// GROUPS

func (r *TestNotesRepository) AddGroup(email string, nameGroup string, pid int) error {
	res, err := r.db.Exec("INSERT INTO test_groups(user_email, name) VALUES ($1, $2)", email, nameGroup)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidData
	}
	return nil
}

func (r *TestNotesRepository) DelGroup(id int, email string) error {
	res, err := r.db.Exec("DELETE FROM test_groups WHERE id = $1 AND user_email = $2", id, email)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidData
	}
	return nil
}

func (r *TestNotesRepository) UpdateGroup(id int, email string, newNameGroup string, pid int) error {
	res, err := r.db.Exec("UPDATE test_groups SET name = $1 WHERE id = $2 AND user_email = $3", newNameGroup, id, email)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidData
	}
	return nil
}

// NOTES

func (r *TestNotesRepository) AddNote(email string, title string, text string, group_id int) error {
	if group_id == 0 {
		_, err := r.db.Exec("INSERT INTO test_notes(user_email, title, text) VALUES ($1, $2, $3)", email, title, text)
		return err
	} else {
		_, err := r.db.Exec(`INSERT INTO test_notes(user_email, title, text, group_id) 
							VALUES ($1, $2, $3, (SELECT id as group_id FROM test_groups WHERE id = $4 AND user_email = $5))`,
			email, title, text, group_id, email)
		return err
	}
}

func (r *TestNotesRepository) DelNote(id int, email string) error {
	res, err := r.db.Exec("DELETE FROM test_notes WHERE id = $1 AND user_email = $2", id, email)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidData
	}
	return nil
}

func (r *TestNotesRepository) UpdateNote(id int, email string, title string, text string, group_id int) error {
	var res sql.Result
	var err error
	if group_id == 0 {
		res, err = r.db.Exec("UPDATE test_notes SET title = $1, text = $2 WHERE id = $3 AND user_email = $4", title, text, id, email)
	} else {
		res, err = r.db.Exec("UPDATE test_notes SET title = $1, text = $2, group_id = $3 WHERE id = $4 AND user_email = $5", title, text, id, group_id, email)
	}
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidData
	}
	return nil
}

func (r *TestNotesRepository) GetNotesList(email string) (model.NoteList, error) {
	// var res *sql.Rows
	// var err error
	// var list model.NoteList

	// if group_id != 0 {
	// 	res, err = r.db.Query(
	// 		"SELECT id, title, group_id FROM test_notes WHERE user_email = $1 AND group_id = $2",
	// 		email,
	// 		group_id,
	// 	)
	// } else {
	// 	res, err = r.db.Query(
	// 		"SELECT id, title, COALESCE(group_id,0) as group_id FROM test_notes WHERE user_email = $1",
	// 		email,
	// 	)
	// }
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return list, nil
	// 	}
	// 	return list, err
	// }
	// defer res.Close()

	// for res.Next() {
	// 	var id int
	// 	var title string
	// 	var group_id int
	// 	if err := res.Scan(&id, &title, &group_id); err != nil {
	// 		return list, err
	// 	}
	// 	list = append(list, struct {
	// 		Id       int    `json:"id"`
	// 		Title    string `json:"title"`
	// 		Group_id int    `json:"group_id"`
	// 	}{
	// 		id,
	// 		title,
	// 		group_id,
	// 	})
	// }

	// if err = res.Err(); err != nil {
	// 	return list, err
	// }
	// return list, nil
	return model.NoteList{}, nil
}

func (r *TestNotesRepository) GetNote(id int, email string) (model.Note, error) {

	var note model.Note

	err := r.db.QueryRow(
		"SELECT id, user_email, title, text, COALESCE(group_id,0) FROM test_notes WHERE id = $1 AND user_email = $2",
		id,
		email,
	).Scan(&note.Id, &note.User_email, &note.Title, &note.Text, &note.Group_id)
	if err != nil {
		return note, err
	}

	return note, nil
}
