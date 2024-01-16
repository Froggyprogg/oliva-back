package main

import (
	"log"
	"oliva-back/pkg/app"
	"oliva-back/pkg/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	application := app.Run(cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
	log.Print("Gracefully stopped")
}
