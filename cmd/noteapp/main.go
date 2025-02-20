package main

import (
	"context"
	"noteapp/internal/server"
	"noteapp/pkg/logger"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func init() {
	logger.SetLevel(6)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer logger.NewLog("main.go", 5, nil, "Stop server", nil)

	server.Start(ctx)
}
