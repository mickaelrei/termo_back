package entities

type GameLetterStatus int8
type GameWordStatus []GameLetterStatus
type GameStatus []GameWordStatus

const (
	// GameLetterStatusCorrect is returned when a letter is in the correct position
	GameLetterStatusCorrect GameLetterStatus = iota

	// GameLetterStatusWrongPosition is returned when a letter is in the wrong position
	GameLetterStatusWrongPosition

	// GameLetterStatusBlack is returned when a letter is not in the word
	GameLetterStatusBlack
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
	WordLength   uint32       `json:"word_length"`
	GameCount    uint32       `json:"game_count"`
	MaxTries     uint32       `json:"max_tries"`
	Attempts     []string     `json:"attempts"`
	GameStatuses []GameStatus `json:"game_statuses"`
}

func (g Game) ToResponse(statuses []GameStatus, maxTries uint32) GameResponse {
	return GameResponse{
		WordLength:   g.GetWordLength(),
		GameCount:    g.GetCount(),
		MaxTries:     maxTries,
		Attempts:     g.Attempts,
		GameStatuses: statuses,
	}
}

func (g Game) GetWordLength() uint32 {
	if len(g.Words) == 0 {
		return 0
	}
	return uint32(len(g.Words[0]))
}

func (g Game) GetCount() uint32 {
	if len(g.Words) == 0 {
		return 0
	}
	return uint32(len(g.Words))
}
