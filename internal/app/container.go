package app

import "github.com/schodevio/trellgo/internal/config"

type Container struct{}

func NewContainer(cfg *config.Config) *Container {
	return &Container{}
}
