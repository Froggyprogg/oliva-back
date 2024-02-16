package main

import (
	"oliva-back/internal/app"
	"oliva-back/internal/config"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
