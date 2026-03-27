package main

import (
	"github.com/schodevio/trellgo/internal/app"
	"github.com/schodevio/trellgo/internal/platform/config"
)

func main() {
	cfg := config.Load()

	server := app.NewServer(cfg)
	server.Start()
}
