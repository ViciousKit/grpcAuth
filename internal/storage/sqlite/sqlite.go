package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sso/internal/domain/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	panic("t")
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	panic("t")

}
func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	panic("t")

}

func (s *Storage) App(ctx context.Context, appId int) (models.App, error) {
	panic("t")

}
