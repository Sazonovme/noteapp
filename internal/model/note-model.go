package model

type Note struct {
	Id         int    `json:"id"`
	User_login string `json:"user_login"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Group_id   int    `json:"group_id"`
}

type Group struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	User_login string `json:"user_login"`
}
