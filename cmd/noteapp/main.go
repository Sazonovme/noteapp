package main

import (
	"encoding/json"
	"noteapp/internal/server"
	"noteapp/pkg/logger"
	"os"

	_ "github.com/lib/pq"
)

func init() {
	logger.SetLevel(6)
}

func main() {
	data, err := os.ReadFile("../../configs/config.json")
	if err != nil {
		logger.NewLog("main - main()", 1, nil, "Filed to read config.json", nil)
		return
	}

	config := server.NewConfig()
	if err = json.Unmarshal(data, config); err != nil {
		logger.NewLog("main - main()", 1, nil, "Filed to unmarshal JSON", nil)
		return
	}

	logger.SetLevel(config.LogLevel)
	server.Start(config)
}
