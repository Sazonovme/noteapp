package server

type configServer struct {
	LogLevel int    `json:"log-level"`
	Port     int    `json:"port"`
	Addr     string `json:"addr"`
	DataBase struct {
		Host     string `json:"host"`
		Username string `username:"username"`
		Dbname   string `json:"dbname"`
		Sslmode  string `json:"sslmode"`
	} `json:"database"`
}

func NewConfig() *configServer {
	return &configServer{}
}
