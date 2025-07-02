package entities

// User maps data from users in the database
type User struct {
	// ID is the database identifier
	ID int64

	// Name is the user's name
	Name string

	// Password is the user's hashed password
	Password string

	// Score tells how many games the user has won
	Score uint32
}

// UserCredentials stores data for an attempt at user registration/login
type UserCredentials struct {
	// Name is the user's name
	Name string `json:"name"`

	// Password is the user's password attempt
	Password string `json:"password"`
}

// UserResponse is used in endpoints to send the minimum required public data
type UserResponse struct {
	// ID is the user's server identifier
	ID int64 `json:"id"`

	// Name is the user's name
	Name string `json:"name"`

	// Score tells how many games the user has won
	Score uint32 `json:"score"`

	// ActiveGame is the user's active game data
	ActiveGame *GameResponse `json:"active_game"`
}

func (u User) ToResponse(game *Game, gameStatuses []GameState, maxGameAttempts *uint32) UserResponse {
	var gameResponse *GameResponse
	if game != nil && maxGameAttempts != nil {
		response := game.ToResponse(gameStatuses, *maxGameAttempts)
		gameResponse = &response
	}

	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Score:      u.Score,
		ActiveGame: gameResponse,
	}
}
