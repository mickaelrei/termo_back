package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/util"
)

type GameRepository interface {
	// StartGame attempts to register a new game in the database for the provided user
	StartGame(ctx context.Context, userID int64, words []string) error

	// RegisterAttempt attempts to register an attempt on the provided game
	RegisterAttempt(ctx context.Context, gameID int64, attempt string, idx uint32, finish bool) error

	// FinishGame marks a game as finished/inactive
	FinishGame(ctx context.Context, gameID int64) error

	// GetUserActiveGame attempts to find the provided user's active game; returns nil if no active game
	GetUserActiveGame(ctx context.Context, userID int64) (*entities.Game, error)
}

type gameRepo struct {
	db *sql.DB
}

func NewGameRepo(db *sql.DB) GameRepository {
	return gameRepo{
		db: db,
	}
}

func (r gameRepo) StartGame(ctx context.Context, userID int64, words []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("[BeginTx] | %v", err)
	}
	defer util.DeferTxRollback(tx)

	// Insert game
	query := `
	INSERT INTO game (id_user) VALUES (?)
	`

	res, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	gameID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("[LastInsertId] | %v", err)
	}

	// Insert words
	var (
		placeholders []string
		args         []any
	)
	for i, word := range words {
		placeholders = append(placeholders, "(?, ?, ?)")
		args = append(args, gameID, word, i)
	}

	queryWord := `
	INSERT INTO game_word (
		id_game,
		word,
		idx
	) VALUES
	` + strings.Join(placeholders, ",\n")

	_, err = tx.ExecContext(ctx, queryWord, args...)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[Commit] | %v", err)
	}

	return nil
}

func (r gameRepo) RegisterAttempt(
	ctx context.Context,
	gameID int64,
	attempt string,
	idx uint32,
	finish bool,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("[BeginTx] | %v", err)
	}
	defer util.DeferTxRollback(tx)

	// Insert the attempt
	queryAttempt := `
	INSERT INTO game_attempt (
		id_game,
		attempt,
		idx
	) VALUES (?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, queryAttempt, gameID, attempt, idx)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	if finish {
		// Finish the game
		queryFinish := `
		UPDATE game
		SET is_active = FALSE
		WHERE id = ?
		`

		_, err = tx.ExecContext(ctx, queryFinish, gameID)
		if err != nil {
			return fmt.Errorf("[ExecContext] | %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[Commit] | %v", err)
	}

	return nil
}

func (r gameRepo) FinishGame(ctx context.Context, gameID int64) error {
	query := `
	UPDATE game
	SET is_active = FALSE
	WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, gameID)
	if err != nil {
		return fmt.Errorf("[ExecContext] | %v", err)
	}

	return nil
}

func (r gameRepo) GetUserActiveGame(ctx context.Context, userID int64) (*entities.Game, error) {
	query := `
	SELECT id
	FROM game
	WHERE id_user = ?
	  AND is_active = TRUE
	`

	var game entities.Game
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&game.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[ExecContext] | %v", err)
	}

	// Get game words
	game.Words, err = r.getGameWords(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("[getGameWords] | %v", err)
	}

	game.Attempts, err = r.getGameAttempts(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("[getGameAttempts] | %v", err)
	}

	return &game, nil
}

func (r gameRepo) getGameWords(ctx context.Context, gameID int64) ([]string, error) {
	query := `
	SELECT word
	FROM game_word
	WHERE id_game = ?
	ORDER BY idx 
	`

	rows, err := r.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("[QueryContext] | %v", err)
	}
	defer util.DeferRowsClose(rows)

	var words []string
	for rows.Next() {
		var word string
		err := rows.Scan(&word)
		if err != nil {
			return nil, fmt.Errorf("[Scan] | %v", err)
		}

		words = append(words, word)
	}

	return words, nil
}

func (r gameRepo) getGameAttempts(ctx context.Context, gameID int64) ([]string, error) {
	query := `
	SELECT attempt
	FROM game_attempt
	WHERE id_game = ?
	ORDER BY idx
	`

	rows, err := r.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("[QueryContext] | %v", err)
	}
	defer util.DeferRowsClose(rows)

	var attempts []string
	for rows.Next() {
		var attempt string
		err := rows.Scan(&attempt)
		if err != nil {
			return nil, fmt.Errorf("[Scan] | %v", err)
		}

		attempts = append(attempts, attempt)
	}

	return attempts, nil
}
