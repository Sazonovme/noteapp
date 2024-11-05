package server

type ConfigServer struct {
	Port          string `json:"port"`
	Db_connection string `json:"db_connection"`
}

func NewConfig() *ConfigServer {
	return &ConfigServer{}
}
