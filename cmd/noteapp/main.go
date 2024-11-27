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

	server.Start(ctx)
}
