package app

import (
	"log/slog"
	grpcApp "sso/internal/app/grpc"
	authService "sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GrpcServer *grpcApp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTtl time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := authService.New(log, storage, storage, storage, tokenTtl)
	grpcApp := grpcApp.New(log, port, authService)

	return &App{
		grpcApp,
	}
}
