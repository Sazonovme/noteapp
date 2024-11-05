package main

import (
	"encoding/json"
	"flag"
	"noteapp/internal/noteapp/server"
	"noteapp/pkg/logger"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
	logLevel   int
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.json", "path to config file (JSON)")
	flag.IntVar(&logLevel, "log-level", 3, "log level for Loggrus")
}

func main() {

	flag.Parse()
	logrus.SetLevel(logrus.Level(logLevel))

	data, err := os.ReadFile("../../" + configPath)
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
