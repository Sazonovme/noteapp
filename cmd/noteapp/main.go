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
		logger.NewLog("main", "main", err, nil, 1, "Can`t read config.json")
	}

	config := server.NewConfig()
	if err = json.Unmarshal(data, config); err != nil {
		logger.NewLog("main", "main", err, nil, 1, "Can`t unmarshal config.json")
	}

	err = server.Start(config)
	if err != nil {
		logger.NewLog("main", "main", err, nil, 1, "Server can`t start")
	}

}
