package entities

// User maps data from users in the database
type User struct {
	// ID is the database identifier
	ID int64

	// Name is the user's name
	Name string

	// Password is the user's hashed password
	Password string
}

// UserCredentials stores data for an attempt at user registration/login
type UserCredentials struct {
	// Name is the user's name
	Name string

	// Password is the user's password attempt
	Password string
}

// UserResponse is used in endpoints to send the minimum required public data
type UserResponse struct {
	// ID is the user's server identifier
	ID int64 `json:"id"`

	// Name is the user's name
	Name string `json:"name"`
}
