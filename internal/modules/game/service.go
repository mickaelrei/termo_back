package game

import (
	"context"
	"termo_back_end/internal/entities"
)

type Service interface {
	StartGame(ctx context.Context, user *entities.User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

func (s service) StartGame(ctx context.Context, user *entities.User) error {
	// TODO: Implement
	return nil
}
