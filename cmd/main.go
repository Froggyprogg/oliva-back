package main

import (
	"oliva-back/pkg/app"
	"oliva-back/pkg/config"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)

}
