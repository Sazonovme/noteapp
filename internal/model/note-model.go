package model

type NoteList []struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Group_id int    `json:"group_id"`
}

type Note struct {
	Id         int    `json:"id"`
	User_email string `json:"user_email"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Group_id   int    `json:"group_id"`
}

type GroupList []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	User_email string `json:"user_email"`
}
