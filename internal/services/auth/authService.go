package authService

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/services/storage"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

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
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.Any("err", err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		log.Error("failed to get app", slog.Any("err", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", slog.Any("err", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
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
	const op = "auth.Admin"
	log := a.log.With(slog.String("op", op))
	log.Info("checking is user admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, int64(userId))
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
