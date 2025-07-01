package rules

import (
	"fmt"
	"termo_back_end/internal/entities"
)

const letterBlank = byte('\n')

// GetGameMaxAttempts returns the maximum number of attempts a user can make for a given word length and game count
func GetGameMaxAttempts(wordLength uint32, gameCount uint32) uint32 {
	// TODO: Incorporate wordLength into the calculation
	return gameCount + 5
}

// CheckGameAttempt checks a word attempt at a game. Returns a GameWordState for each game word, containing the
// GameLetterState for each letter in the word.
//
// Note: The input attempt and the game words are expected to have the same length and to be trimmed, lowercased and
// cleaned (no diacritics)
func CheckGameAttempt(game entities.Game, attempt string) []entities.GameWordState {
	gameStatus := make([]entities.GameWordState, len(game.Words))

	for j, word := range game.Words {
		wordStatus := make([]entities.GameLetterState, len(attempt))

		wordCopy := make([]byte, len(word))
		attemptCopy := make([]byte, len(word))
		copy(wordCopy, word)
		copy(attemptCopy, attempt)

		// First, mark the letters in the correct position
		for i := range attempt {
			if attempt[i] == word[i] {
				wordStatus[i] = entities.GameLetterStateCorrect
				wordCopy[i] = letterBlank
				attemptCopy[i] = letterBlank
			} else {
				wordStatus[i] = -1
			}
		}

		// Then, mark correct letters in the wrong position
		for i := range attempt {
			if wordStatus[i] != -1 || attemptCopy[i] == letterBlank {
				continue
			}

			idx := index(wordCopy, attempt[i])
			if idx != -1 {
				wordStatus[i] = entities.GameLetterStateWrongPosition
				wordCopy[idx] = letterBlank
			}
		}

		// Finally, mark incorrect letters as black
		for i := range attempt {
			if wordStatus[i] == -1 {
				wordStatus[i] = entities.GameLetterStateWrong
			}
		}

		gameStatus[j] = wordStatus
	}

	return gameStatus
}

func IsGameWon(game entities.Game, currentAttempt string) bool {
	// Every word must have at least one match
	for _, word := range game.Words {
		match := false
		// Check if any previous attempt is wrong
		for _, prevAttempt := range game.Attempts {
			fmt.Printf("comparing \"%s\" and \"%s\"\n", prevAttempt, word)
			if prevAttempt == word {
				fmt.Println("found match")
				match = true
				break
			}
		}

		// Check if the current attempt is wrong
		fmt.Printf("comparing \"%s\" and \"%s\"\n", currentAttempt, word)
		if currentAttempt == word {
			fmt.Println("found match curr")
			match = true
		}

		if !match {
			return false
		}
	}
	return true
}

func index(s []byte, b byte) int {
	for i, v := range s {
		if v == b {
			return i
		}
	}
	return -1
}
