package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"noteapp/internal/model"
	"noteapp/pkg/logger"
	"strconv"
)

var (
	ErrInvalidData    = errors.New("incorrect data")
	ErrNoDataChanges  = errors.New("no data found to change")
	ErrFiledToConvert = errors.New("filed to convert string to int")
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

func (r *NotesRepository) UpdateGroup(id int, email string, newNameGroup string, pid int) error {

	var res sql.Result
	var err error

	if newNameGroup != "" && pid == -1 {
		res, err = r.db.Exec("UPDATE groups SET name = $1 WHERE id = $2 AND user_email = $3", newNameGroup, id, email)
	} else if newNameGroup == "" && pid != -1 {
		res, err = r.db.Exec("UPDATE groups SET pid = $1 WHERE id = $2 AND user_email = $3", pid, id, email)
	} else if newNameGroup != "" && pid != -1 {
		res, err = r.db.Exec("UPDATE groups SET name = $1, pid = $2 WHERE id = $3 AND user_email = $4", newNameGroup, pid, id, email)
	} else {
		return ErrNoDataChanges
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

// NOTES

func (r *NotesRepository) AddNote(email string, title string, group_id int) error {
	if group_id == -1 {
		_, err := r.db.Exec("INSERT INTO notes(user_email, title) VALUES ($1, $2)", email, title)
		return err
	} else {
		_, err := r.db.Exec(`INSERT INTO notes(user_email, title, group_id) 
							VALUES ($1, $2, (SELECT id as group_id FROM groups WHERE id = $3 AND user_email = $4))`,
			email, title, group_id, email)
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

func (r *NotesRepository) UpdateNote(data map[string]string) error {

	sqlRequest, params, err := getRequestAndParams(data)
	fmt.Println(sqlRequest)
	fmt.Println(params...)
	if err != nil {
		return err
	}

	if sqlRequest == "" {
		return ErrInvalidData
	}

	res, err := r.db.Exec(sqlRequest, params...)
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
			WHERE user_email = $1 AND pid = 0 

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
		); err != nil {
			return model.NoteList{}, err
		}

		if resRow.group_id == 0 {
			notes = append(notes, model.NoteElement{
				Id:    resRow.notes_id,
				Title: resRow.notes_title,
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
				})
			}
		}
	}

	if curGrp.Groups != nil {
		*groupsLink = append(*groupsLink, curGrp)
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

// HELPER

func getRequestAndParams(data map[string]string) (string, []interface{}, error) {

	params := []interface{}{}
	sqlRequest := ""

	string_id := data["id"]
	id, err := strconv.Atoi(string_id)
	if err != nil {
		logger.NewLog("repo - getRequestAndParams()", 2, err, "Filed to convert string to int", "string = "+string_id)
		return "", nil, ErrFiledToConvert
	}
	email := data["email"]

	text, textOK := data["text"]
	title, titleOK := data["title"]
	group_id_string, group_id_stringOK := data["group_id"]

	count := 1

	if titleOK {
		sqlRequest += " title = $" + strconv.Itoa(count) + " "
		params = append(params, title)
		count++
	}

	if textOK {
		if sqlRequest != "" {
			sqlRequest += ","
		}
		sqlRequest += " text = $" + strconv.Itoa(count) + " "
		params = append(params, text)
		count++
	}

	if group_id_stringOK {

		group_id, err := strconv.Atoi(group_id_string)
		if err != nil {
			logger.NewLog("repo - getRequestAndParams()", 2, err, "Filed to convert string to int", "string = "+group_id_string)
			return "", nil, ErrFiledToConvert
		}

		if sqlRequest != "" {
			sqlRequest += ","
		}
		sqlRequest += " group_id = (SELECT id as group_id FROM groups WHERE id = $" + strconv.Itoa(count) + " AND user_email = $" + strconv.Itoa(count+1) + ") "
		params = append(params, group_id)
		params = append(params, email)
		count += 2
	}

	if len(params) < 1 {
		return "", nil, ErrNoDataChanges
	}

	sqlRequest = "UPDATE notes SET" + sqlRequest + "WHERE id = $" + strconv.Itoa(count) + " AND user_email = $" + strconv.Itoa(count+1)
	params = append(params, id)
	params = append(params, email)

	return sqlRequest, params, nil
}
