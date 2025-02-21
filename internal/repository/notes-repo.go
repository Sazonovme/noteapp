package repository

import (
	"database/sql"
	"errors"
	"noteapp/internal/model"
)

var (
	ErrInvalidData = errors.New("incorrect data")
)

type NotesRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *NotesRepository {
	return &NotesRepository{
		db: db,
	}
}

// GROUPS

func (r *NotesRepository) AddGroup(email string, nameGroup string, pid int) error {
	var res sql.Result
	var err error

	if pid == 0 {
		res, err = r.db.Exec("INSERT INTO groups(user_email, name, pid) VALUES ($1, $2, $3)", email, nameGroup, nil)
	} else {
		res, err = r.db.Exec("INSERT INTO groups(user_email, name, pid) VALUES ($1, $2, $3)", email, nameGroup, pid)
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

func (r *NotesRepository) DelGroup(id int, email string) error {
	res, err := r.db.Exec("DELETE FROM groups WHERE id = $1 AND user_email = $2", id, email)
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

func (r *NotesRepository) UpdateGroup(id int, email string, newNameGroup string) error {
	res, err := r.db.Exec("UPDATE groups SET name = $1 WHERE id = $2 AND user_email = $3", newNameGroup, id, email)
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

func (r *NotesRepository) AddNote(email string, title string, text string, group_id int) error {
	if group_id == 0 {
		_, err := r.db.Exec("INSERT INTO notes(user_email, title, text) VALUES ($1, $2, $3)", email, title, text)
		return err
	} else {
		_, err := r.db.Exec(`INSERT INTO notes(user_email, title, text, group_id) 
							VALUES ($1, $2, $3, (SELECT id as group_id FROM groups WHERE id = $4 AND user_email = $5))`,
			email, title, text, group_id, email)
		return err
	}
}

func (r *NotesRepository) DelNote(id int, email string) error {
	res, err := r.db.Exec("DELETE FROM notes WHERE id = $1 AND user_email = $2", id, email)
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

func (r *NotesRepository) UpdateNote(id int, email string, title string, text string, group_id int) error {
	var res sql.Result
	var err error
	if group_id == 0 {
		res, err = r.db.Exec("UPDATE notes SET title = $1, text = $2 WHERE id = $3 AND user_email = $4", title, text, id, email)
	} else {
		res, err = r.db.Exec("UPDATE notes SET title = $1, text = $2, group_id = $3 WHERE id = $4 AND user_email = $5", title, text, id, group_id, email)
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

func (r *NotesRepository) GetNotesList(email string) (model.NoteList, error) {
	var res *sql.Rows

	res, err := r.db.Query(
		`WITH RECURSIVE r AS (
			SELECT id, pid, name, 1 AS level
			FROM groups
			WHERE user_email = $1 AND pid IS NULL 

			UNION

			SELECT groups.id, groups.pid, groups.name, r.level + 1 AS level
			FROM groups
				JOIN r
					ON groups.pid = r.id
		)
		SELECT
			COALESCE(r.id, 0) AS group_id,
			COALESCE(r.name, '') AS group_name,
			COALESCE(r.pid, 0) AS group_pid,
			COALESCE(r.level, 1) AS group_level,
			COALESCE(notes.id, 0) AS notes_id,
			COALESCE(notes.title, '') AS notes_title,
			COALESCE(notes.text, '') AS notes_text
		FROM r 
			FULL OUTER JOIN notes
				ON notes.group_id = r.id
		ORDER BY group_level ASC, group_pid ASC, group_id ASC;`,
		email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.NoteList{}, nil
		}
		return model.NoteList{}, err
	}
	defer res.Close()

	resRow := struct {
		group_id    int
		group_name  string
		group_pid   int
		group_level int
		notes_id    int
		notes_title string
		notes_text  string
	}{}

	gPid := 0
	gId := -1

	curGrp := model.GroupElement{}
	groups := []model.GroupElement{}
	groupsLink := &groups
	notes := []model.NoteElement{}

	mGroups := map[int]model.GroupElement{}
	defer clear(mGroups)

	for res.Next() {

		if err := res.Scan(
			&resRow.group_id,
			&resRow.group_name,
			&resRow.group_pid,
			&resRow.group_level,
			&resRow.notes_id,
			&resRow.notes_title,
			&resRow.notes_text,
		); err != nil {
			return model.NoteList{}, err
		}

		if resRow.group_id == 0 {
			notes = append(notes, model.NoteElement{
				Id:    resRow.notes_id,
				Title: resRow.notes_title,
				Text:  resRow.notes_text,
			})
			continue
		} else {
			if gId != resRow.group_id {

				if gId != -1 {
					*groupsLink = append(*groupsLink, curGrp)

					mGroups[curGrp.Id] = curGrp
				}

				if gPid != resRow.group_pid {
					groupElement, ok := mGroups[resRow.group_pid]
					if !ok {
						return model.NoteList{}, errors.New("not found group in map")
					}
					groupsLink = groupElement.Groups
				}

				curGrp = model.GroupElement{}
				curGrp.Id = resRow.group_id
				curGrp.Name = resRow.group_name
				curGrp.Groups = &[]model.GroupElement{}
				curGrp.Notes = []model.NoteElement{}

				gId = resRow.group_id
				gPid = resRow.group_pid
			}

			if resRow.notes_id != 0 {
				curGrp.Notes = append(curGrp.Notes, model.NoteElement{
					Id:    resRow.notes_id,
					Title: resRow.notes_title,
					Text:  resRow.notes_text,
				})
			}
		}
	}

	if err = res.Err(); err != nil {
		return model.NoteList{}, err
	}

	return model.NoteList{
		Groups: groups,
		Notes:  notes,
	}, nil
}

func (r *NotesRepository) GetNote(id int, email string) (model.Note, error) {

	var note model.Note

	err := r.db.QueryRow(
		"SELECT id, user_email, title, text, COALESCE(group_id,0) FROM notes WHERE id = $1 AND user_email = $2",
		id,
		email,
	).Scan(&note.Id, &note.User_email, &note.Title, &note.Text, &note.Group_id)
	if err != nil {
		return note, err
	}

	return note, nil
}
