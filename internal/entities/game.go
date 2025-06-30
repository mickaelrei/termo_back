package entities

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
