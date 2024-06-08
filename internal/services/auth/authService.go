package authService

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	userSaver    UserSaver
	userProvider UserProvider
	log          *slog.Logger
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int) (token string, err error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))
	log.Info("login")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		log.Error("failed to get user", slog.Any("err", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op))
	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed gen pass hash", slog.Any("err", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.Any("err", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)

	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int) (bool, error) {
	panic("t")
}
