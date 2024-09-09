package main

import (
	"fmt"
	"log/slog"
	"tenders/internal/app"
	"tenders/internal/config"
)

const envFilePath = "/usr/local/.env"

func main() {
	cfg, err := config.Parse(envFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	slog.Info("Server is starting")

	app.Run(cfg)
}
