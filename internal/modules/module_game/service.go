package module_game

import (
	"context"
	"errors"
	"fmt"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
)

type Service interface {
	// StartGame attempts to start a game for the provided user with the given configs
	StartGame(
		ctx context.Context,
		user *entities.User,
		wordLength uint32,
		gameCount uint32,
	) (status_codes.GameStart, error)

	// AttemptGame attempts to register an attempt on the current game of the provided user
	AttemptGame(ctx context.Context, user *entities.User, attempt string) (status_codes.GameAttempt, error)

	// GetUserActiveGame attempts to find the provided user's active game; returns nil if no active game
	GetUserActiveGame(ctx context.Context, user *entities.User) (*entities.Game, error)
}

type service struct {
	wordMap util.WordMap
	repo    Repository
}

func NewService(words []string, repo Repository) Service {
	return service{
		wordMap: util.WordMapFromList(words),
		repo:    repo,
	}
}

func (s service) StartGame(
	ctx context.Context,
	user *entities.User,
	wordLength uint32,
	gameCount uint32,
) (status_codes.GameStart, error) {
	// Check if the user is already in a game
	game, err := s.repo.GetUserActiveGame(ctx, user.ID)
	if err != nil {
		return -1, fmt.Errorf("[GetUserActiveGame] | %v", err)
	}

	if game != nil {
		return status_codes.GameStartActiveGame, nil
	}

	// Ensure valid configs
	// TODO: Better way of handling this; hardcoded for now
	if wordLength < 3 || wordLength > 22 {
		return status_codes.GameStartInvalidWordLength, nil
	}
	if gameCount == 0 || gameCount > 20 {
		return status_codes.GameStartInvalidCount, nil
	}

	// Choose words randomly
	words, err := s.wordMap.ChooseRandom(wordLength, gameCount)
	if err != nil {
		if errors.Is(err, util.ErrInvalidSize) {
			return status_codes.GameStartInvalidWordLength, nil
		}
		if errors.Is(err, util.ErrNotEnoughWords) {
			return status_codes.GameStartInvalidCount, nil
		}
		return -1, fmt.Errorf("[ChooseRandom] | %v", err)
	}

	// Register game in the database
	err = s.repo.StartGame(ctx, user.ID, words)
	if err != nil {
		return -1, fmt.Errorf("[StartGame] | %v", err)
	}

	return status_codes.GameStartSuccess, nil
}

func (s service) AttemptGame(
	ctx context.Context,
	user *entities.User,
	attempt string,
) (status_codes.GameAttempt, error) {
	// Ensure the user is already in a game
	game, err := s.GetUserActiveGame(ctx, user)
	if err != nil {
		return -1, fmt.Errorf("[GetUserActiveGame] | %v", err)
	}

	if game == nil {
		return status_codes.GameAttemptNoActiveGame, nil
	}

	// Ensure the attempt is valid
	if len(attempt) != 5 {
		return status_codes.GameAttemptInvalid, nil
	}

	// TODO: Check what's right and what's wrong
	var finish bool

	// Register attempt in database
	err = s.repo.RegisterAttempt(ctx, game.ID, attempt, finish)
	if err != nil {
		return -1, fmt.Errorf("[RegisterAttempt] | %v", err)
	}

	return status_codes.GameAttemptSuccess, nil
}

func (s service) GetUserActiveGame(
	ctx context.Context,
	user *entities.User,
) (*entities.Game, error) {
	return s.repo.GetUserActiveGame(ctx, user.ID)
}
