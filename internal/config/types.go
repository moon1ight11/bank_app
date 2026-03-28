package config

import "time"

type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Logger      Logger         `mapstructure:"logger"`
}

type Logger struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"filepath"`
	MaxSize  int    `mapstructure:"maxsize"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Name          string `mapstructure:"name"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	MigrationsDir string `mapstructure:"migrationsDir"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}
