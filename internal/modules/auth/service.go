package auth

import (
	"context"
	"termo_back_end/internal/entities"
)

type Service interface {
	RegisterUser(ctx context.Context, user *entities.User) error
	LoginUser(ctx context.Context, user *entities.User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) RegisterUser(ctx context.Context, user *entities.User) error {
	// TODO: Check if user already exists

	return s.repo.RegisterUser(ctx, user)
}

func (s service) LoginUser(ctx context.Context, user *entities.User) error {
	// TODO: Check if user exists

	// TODO: Check if password matches

	return nil
}
