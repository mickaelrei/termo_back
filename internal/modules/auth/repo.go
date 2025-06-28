package auth

import (
	"context"
	"database/sql"
	"fmt"
	"termo_back_end/internal/entities"
)

type Repository interface {
	RegisterUser(ctx context.Context, user *entities.User) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return repo{
		db: db,
	}
}

func (r repo) RegisterUser(ctx context.Context, user *entities.User) error {
	query := `
	INSERT INTO user (
		name,
		password,
		salt
	) VALUES (?, ?, ?)
	`

	// TODO: Generate password salt (or receive it from the service)
	var salt string

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Password, salt)
	if err != nil {
		return fmt.Errorf("error in [ExecContext]: %v", err)
	}

	return nil
}
