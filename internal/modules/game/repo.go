package game

import (
	"context"
	"database/sql"
	"termo_back_end/internal/entities"
)

type Repository interface {
	StartGame(ctx context.Context, user *entities.User) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return repo{
		db: db,
	}
}

func (r repo) StartGame(ctx context.Context, user *entities.User) error {
	// TODO: Implement
	return nil
}
