package service

import (
	"context"
	"errors"
	"fmt"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/repo"
	"termo_back_end/internal/rules"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
)

type GameAttemptData struct {
	Status    status_codes.GameAttempt
	GameState []entities.GameWordState
	Words     []string
	Won       bool
}

type GameService interface {
	// StartGame attempts to start a game for the provided user with the given configs
	StartGame(
		ctx context.Context,
		user *entities.User,
		wordLength uint32,
		wordCount uint32,
	) (status_codes.GameStart, error)

	// AttemptGame attempts to register an attempt on the current game of the provided user
	AttemptGame(
		ctx context.Context,
		user *entities.User,
		attempt string,
	) (*GameAttemptData, error)

	// GetUserActiveGame attempts to find the provided user's active game; returns nil if no active game
	GetUserActiveGame(
		ctx context.Context,
		user *entities.User,
	) (*entities.Game, []entities.GameState, error)
}

type gameService struct {
	wordMap  util.WordMap
	repo     repo.GameRepository
	userRepo repo.UserRepository
}

func NewGameService(words []string, repo repo.GameRepository, userRepo repo.UserRepository) GameService {
	return gameService{
		wordMap:  util.WordMapFromList(words),
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s gameService) StartGame(
	ctx context.Context,
	user *entities.User,
	wordLength uint32,
	wordCount uint32,
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
	if wordCount == 0 || wordCount > 20 {
		return status_codes.GameStartInvalidCount, nil
	}

	// Choose words randomly
	words, err := s.wordMap.ChooseRandom(wordLength, wordCount)
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

func (s gameService) AttemptGame(
	ctx context.Context,
	user *entities.User,
	attempt string,
) (*GameAttemptData, error) {
	// Clean attempt
	attempt = s.wordMap.CleanWord(attempt)

	// Ensure the user is already in a game
	game, err := s.repo.GetUserActiveGame(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("[GetUserActiveGame] | %v", err)
	}

	if game == nil {
		return &GameAttemptData{
			Status: status_codes.GameAttemptNoActiveGame,
		}, nil
	}

	// Ensure the attempt is valid
	if uint32(len(attempt)) != game.GetWordLength() {
		return &GameAttemptData{
			Status: status_codes.GameAttemptInvalid,
		}, nil
	}

	// Check what's right and what's wrong
	gameState := rules.CheckGameAttempt(*game, attempt)
	currentAttempts := uint32(len(game.Attempts))
	maxAttempts := rules.GetGameMaxAttempts(game.GetWordLength(), game.GetWordCount())

	// Register attempt in database
	err = s.repo.RegisterAttempt(ctx, game.ID, attempt, currentAttempts, currentAttempts >= maxAttempts-1)
	if err != nil {
		return nil, fmt.Errorf("[RegisterAttempt] | %v", err)
	}

	// If all words are correct, increment the user's score
	won := rules.IsGameWon(*game, attempt)
	if won {
		err = s.userRepo.IncrementScore(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("[IncrementScore] | %v", err)
		}

		err = s.repo.FinishGame(ctx, game.ID)
		if err != nil {
			return nil, fmt.Errorf("[FinishGame] | %v", err)
		}
	}

	var words []string
	if currentAttempts >= maxAttempts-1 || won {
		// Player either won or lost; show actual words
		for _, word := range game.Words {
			original, ok := s.wordMap.GetOriginalWord(word)

			if ok {
				words = append(words, original)
			} else {
				words = append(words, word)
			}
		}
	}

	return &GameAttemptData{
		Status:    status_codes.GameAttemptSuccess,
		GameState: gameState,
		Words:     words,
		Won:       won,
	}, nil
}

func (s gameService) GetUserActiveGame(
	ctx context.Context,
	user *entities.User,
) (*entities.Game, []entities.GameState, error) {
	// Get game
	game, err := s.repo.GetUserActiveGame(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("[GetUserActiveGame] | %v", err)
	}

	if game == nil {
		return nil, nil, nil
	}

	// Get status for each attempt
	statuses := make([]entities.GameState, len(game.Attempts))
	for i, attempt := range game.Attempts {
		statuses[i] = rules.CheckGameAttempt(*game, attempt)
	}

	return game, statuses, nil
}
