// @title           TrellGo API
// @version         1.0
// @description     REST API for a Trello-like app.
// @host            localhost:3000
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter: Bearer <token>
package main

import (
	_ "github.com/schodevio/trellgo/docs"
	"github.com/schodevio/trellgo/internal/app"
	"github.com/schodevio/trellgo/internal/platform/config"
)

func main() {
	cfg := config.Load()

	server := app.NewServer(cfg)
	server.Start()
}
