package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// загрузка конфигурации
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}

// дефолтные значения полей конфигурации
func setDefaults() {
	viper.SetDefault("environment", "development")
}
