package entities

type GameLetterState int8
type GameWordState []GameLetterState
type GameState []GameWordState

const (
	// GameLetterStateCorrect is returned when a letter is in the correct position
	GameLetterStateCorrect GameLetterState = iota

	// GameLetterStateWrongPosition is returned when a letter is in the wrong position
	GameLetterStateWrongPosition

	// GameLetterStateWrong is returned when a letter is not in the word
	GameLetterStateWrong
)

// Game maps data from games in the database
type Game struct {
	// ID is the database identifier
	ID int64

	// Words is a list containing all the game's chosen words
	Words []string

	// Attempts is a list containing all the user attempts on this game
	Attempts []string

	// IsActive tells whether this game is active
	IsActive bool
}

// GameResponse is used in endpoints to send the minimum required public data
type GameResponse struct {
	WordLength  uint32      `json:"word_length"`
	WordCount   uint32      `json:"word_count"`
	MaxAttempts uint32      `json:"max_attempts"`
	Attempts    []string    `json:"attempts"`
	GameStates  []GameState `json:"game_states"`
}

func (g Game) ToResponse(states []GameState, maxAttempts uint32) GameResponse {
	return GameResponse{
		WordLength:  g.GetWordLength(),
		WordCount:   g.GetWordCount(),
		MaxAttempts: maxAttempts,
		Attempts:    g.Attempts,
		GameStates:  states,
	}
}

func (g Game) GetWordLength() uint32 {
	if len(g.Words) == 0 {
		return 0
	}
	return uint32(len(g.Words[0]))
}

func (g Game) GetWordCount() uint32 {
	return uint32(len(g.Words))
}
