package models

// Config contains different types of configurations
type Config struct {
	ServerConfig   ServerConfig   `yaml:"server"`
	PostgresConfig PostgresConfig `yaml:"postgres"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}
