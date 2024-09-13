package main

import (
	"fmt"
	"log/slog"
	"tenders/internal/app"
	"tenders/internal/config"
)

const envFilePath = "/usr/local/.env"

//const envFilePath = "../.env"

func main() {
	cfg, err := config.Parse(envFilePath)
	config.InitDbConnectionString(cfg)

	if err != nil {
		fmt.Println(err)
		return
	}

	slog.Info("Server is starting")

	app.Run(cfg)
}
