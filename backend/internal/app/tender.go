package app

import (
	"tenders/internal/config"
	"tenders/internal/handler"
)

func Run(cfg *config.Config) {

	//ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	handler.NewRouter(cfg)

}
