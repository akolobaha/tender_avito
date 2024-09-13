package app

import (
	"tenders/internal/config"
	"tenders/internal/handler"
)

func Run(cfg *config.Config) {
	handler.NewRouter(cfg)
}
