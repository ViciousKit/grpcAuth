package app

import (
	"log/slog"
	grpcApp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GrpcServer *grpcApp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTtl time.Duration) *App {
	grpcApp := grpcApp.New(log, port)

	return &App{
		grpcApp,
	}
}
