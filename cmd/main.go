package main

import (
	"oliva-back/internal/app"
	"oliva-back/pkg/config"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)

}
