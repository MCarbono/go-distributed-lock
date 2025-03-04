package config

type Config struct {
	Database              RelationalDatabase
	NonRelationalDatabase NonRelationalDatabase
}

type RelationalDatabase struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type NonRelationalDatabase struct {
	Host string
	Port string
}

// func LoadEnv(env string) Config {

// }
