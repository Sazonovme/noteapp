package main

import (
	"encoding/json"
	"noteapp/internal/noteapp/server"
	"noteapp/pkg/logger"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	data, err := os.ReadFile("../../config.json")
	if err != nil {
		logger.NewLog("main.go - main() - os.ReadFile()", 1, nil, "Filed to read config.json", nil)
		return
	}
	config := server.NewConfig()
	if err = json.Unmarshal(data, config); err != nil {
		logger.NewLog("main.go - main() - server.NewConfig()", 1, nil, "Filed to unmarshal JSON", nil)
		return
	}
	server.Start(config)
}
