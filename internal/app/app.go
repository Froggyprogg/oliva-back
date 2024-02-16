package app

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"oliva-back/internal/config"
	user "oliva-back/internal/grpc/auth"
	"oliva-back/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	gRPCServer *grpc.Server
	database   *sql.DB
}

func Run(cfg *config.Config) *App {
	database := storage.NewDatabaseConnection(cfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %d", err)
		return nil
	}

	gRPCServer := grpc.NewServer()
	reflection.Register(gRPCServer)
	user.Register(gRPCServer, database)
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
