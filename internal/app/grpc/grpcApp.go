package grpcApp

import (
	"fmt"
	"log/slog"
	"net"
	"sso/internal/grpc/auth"
	authService "sso/internal/services/auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	grpcServer := grpc.NewServer()

	auth.RegisterServerApi(grpcServer, authService.New())

	return &App{
		log,
		grpcServer,
		port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("started grpcs server", slog.String("address", listener.Addr().String()))

	if err := a.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcApp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping grpcs server", slog.Int("port", a.port))

	a.grpcServer.GracefulStop()
}
