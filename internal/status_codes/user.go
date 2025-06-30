package status_codes

type UserRegister int64
type UserLogin int64
type UserUpdateName int64
type UserUpdatePassword int64

const (
	UserRegisterSuccess UserRegister = iota
	UserRegisterInvalidName
	UserRegisterInvalidPassword
	UserRegisterInvalidAlreadyRegistered
)

const (
	UserLoginSuccess UserLogin = iota
	UserLoginNotFound
	UserLoginWrongPassword
)

const (
	UserUpdateNameSuccess UserUpdateName = iota
	UserUpdateNameInvalid
)

const (
	UserUpdatePasswordSuccess UserUpdatePassword = iota
	UserUpdatePasswordInvalid
)

func (c UserRegister) String() string {
	switch c {
	case UserRegisterSuccess:
		return "SUCCESS"
	case UserRegisterInvalidName:
		return "INVALID_NAME"
	case UserRegisterInvalidPassword:
		return "INVALID_PASSWORD"
	case UserRegisterInvalidAlreadyRegistered:
		return "ALREADY_REGISTERED"
	default:
		return "UNKNOWN"
	}
}

func (c UserLogin) String() string {
	switch c {
	case UserLoginSuccess:
		return "SUCCESS"
	case UserLoginNotFound:
		return "NOT_FOUND"
	case UserLoginWrongPassword:
		return "WRONG_PASSWORD"
	default:
		return "UNKNOWN"
	}
}

func (c UserUpdateName) String() string {
	switch c {
	case UserUpdateNameSuccess:
		return "SUCCESS"
	case UserUpdateNameInvalid:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}

func (c UserUpdatePassword) String() string {
	switch c {
	case UserUpdatePasswordSuccess:
		return "SUCCESS"
	case UserUpdatePasswordInvalid:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}
