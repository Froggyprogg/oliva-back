package app

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"oliva-back/pkg/config"
	user "oliva-back/pkg/handlers"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	gRPCServer *grpc.Server
}

func Run(cfg *config.Config) *App {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %d", err)
		return nil
	}

	gRPCServer := grpc.NewServer()
	reflection.Register(gRPCServer)
	user.Register(gRPCServer)
	log.Printf("server listening at :%d", cfg.GRPC.Port)

	if err = gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %d", err)
		return nil
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")

	gRPCServer.GracefulStop()

	log.Println("Server gracefully stopped")

	return &App{gRPCServer: gRPCServer}
}
