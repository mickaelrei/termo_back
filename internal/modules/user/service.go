package user

import (
	"context"
	"termo_back_end/internal/entities"
)

type Service interface {
	UpdateName(ctx context.Context, user *entities.User) error
	UpdatePassword(ctx context.Context, user *entities.User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) UpdateName(ctx context.Context, user *entities.User) error {
	// TODO: Implement
	return nil
}

func (s service) UpdatePassword(ctx context.Context, user *entities.User) error {
	// TODO: Implement
	return nil
}
