package status_codes

type GameStart int64
type GameAttempt int64

const (
	GameStartSuccess GameStart = iota
	GameStartActiveGame
	GameStartInvalidWordLength
	GameStartInvalidCount
)

const (
	GameAttemptSuccess GameAttempt = iota
	GameAttemptNoActiveGame
	GameAttemptInvalid
)

func (c GameStart) String() string {
	switch c {
	case GameStartSuccess:
		return "SUCCESS"
	case GameStartActiveGame:
		return "ALREADY_IN_PROGRESS"
	case GameStartInvalidWordLength:
		return "INVALID_WORD_LENGTH"
	case GameStartInvalidCount:
		return "INVALID_COUNT"
	default:
		return "UNKNOWN"
	}
}

func (c GameAttempt) String() string {
	switch c {
	case GameAttemptSuccess:
		return "SUCCESS"
	case GameAttemptNoActiveGame:
		return "NO_ACTIVE_GAME"
	case GameAttemptInvalid:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}
