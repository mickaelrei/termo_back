package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"termo_back_end/internal/entities"
)

type UserRepository interface {
	// RegisterUser inserts a user into the database with given credentials; returns it if succeeded
	//
	// Password is expected to be already hashed; will be inserted as is
	RegisterUser(ctx context.Context, credentials entities.UserCredentials) (*entities.User, error)

	// GetUserByID attempts to find a user with the provided ID; returns nil if not found
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)

	// GetUserByName attempts to find a user with the provided name; returns nil if not found
	//
	// The name must be an exact match, meaning upper/lowercase letters won't match
	GetUserByName(ctx context.Context, name string) (*entities.User, error)

	// UpdateName updates a user's name, given their ID
	UpdateName(ctx context.Context, userID int64, name string) error

	// UpdatePassword updates a user's password, given their ID
	//
	// Password is expected to be already hashed; will be inserted as is
	UpdatePassword(ctx context.Context, userID int64, password string) error

	// IncrementScore increments a user's score, given their ID
	IncrementScore(ctx context.Context, userID int64) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return userRepo{
		db: db,
	}
}

func (r userRepo) RegisterUser(
	ctx context.Context,
	credentials entities.UserCredentials,
) (*entities.User, error) {
	query := `
	INSERT INTO user (
		name,
		password
	) VALUES (?, ?)
	`

	res, err := r.db.ExecContext(ctx, query, credentials.Name, credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("[ExecContext] | %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("[LastInsertId] | %v", err)
		return nil, fmt.Errorf("[LastInsertId] | %v", err)
	}

	return r.GetUserByID(ctx, id)
}

func (r userRepo) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	query := `
	SELECT id,
	       name,
	       password,
	       score
	FROM user
	WHERE id = ?
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Score,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("[QueryRowContext] | %v", err)
		return nil, fmt.Errorf("[QueryRowContext] | %v", err)
	}

	return &user, nil
}

func (r userRepo) GetUserByName(ctx context.Context, name string) (*entities.User, error) {
	query := `
	SELECT id,
	       name,
	       password,
	       score
	FROM user
	WHERE name = ?
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Score,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("[QueryRowContext] | %v", err)
		return nil, fmt.Errorf("[QueryRowContext] | %v", err)
	}

	return &user, nil
}

func (r userRepo) UpdateName(ctx context.Context, userID int64, name string) error {
	query := `
	UPDATE user
	SET name = ?
	WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, name, userID)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	return nil
}

func (r userRepo) UpdatePassword(ctx context.Context, userID int64, password string) error {
	query := `
	UPDATE user
	SET password = ?
	WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, password, userID)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	return nil
}

func (r userRepo) IncrementScore(ctx context.Context, userID int64) error {
	query := `
	UPDATE user
	SET score = score + 1
	WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	return nil
}
