package entities

type database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
}

type auth struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// Config is a struct used for loading the config.json file with all project configurations
type Config struct {
	Database database `json:"db"`

	Auth auth `json:"auth"`
}
