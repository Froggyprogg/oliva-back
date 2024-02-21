package main

import (
	"log"
	"oliva-back/internal/app"
	"oliva-back/internal/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	cfg := config.NewConfig()

	app.Run(cfg)
}
