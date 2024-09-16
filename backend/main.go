package main

import (
	"log/slog"
	"tenders/internal/app"
	"tenders/internal/config"
)

// Докер
const envFilePath = "/usr/local/.env"

// Локальный запуск
//const envFilePath = "../.env"

func main() {
	cfg, err := config.Parse(envFilePath)
	config.InitDbConnectionString(cfg)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Server is starting")

	app.Run(cfg)
}
