package user

import (
	"context"
	"database/sql"
	"fmt"
	"termo_back_end/internal/entities"
)

type Repository interface {
	UpdateName(ctx context.Context, user *entities.User) error
	UpdatePassword(ctx context.Context, user *entities.User) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return repo{
		db: db,
	}
}

func (r repo) UpdateName(ctx context.Context, user *entities.User) error {
	query := `
	UPDATE user
	SET name = ?
	WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.ID)
	if err != nil {
		return fmt.Errorf("error in [ExecContext]: %v", err)
	}

	return nil
}

func (r repo) UpdatePassword(ctx context.Context, user *entities.User) error {
	// TODO: Implement
	return nil
}
