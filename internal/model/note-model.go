package model

type NoteElement struct {
	Id    int    `json:"note_id"`
	Title string `json:"note_title"`
	Text  string `json:"note_text"`
}

type GroupElement struct {
	Id     int             `json:"group_id"`
	Name   string          `json:"group_name"`
	Groups *[]GroupElement `json:"groups"`
	Notes  []NoteElement   `json:"notes"`
}

type NoteList struct {
	Groups []GroupElement `json:"groups"`
	Notes  []NoteElement  `json:"notes"`
}

type Note struct {
	Id         int    `json:"id"`
	User_email string `json:"user_email"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Group_id   int    `json:"group_id"`
}

type Group struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	User_email string `json:"user_email"`
}
