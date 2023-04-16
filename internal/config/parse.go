package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/razielsd/chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("error during load config: %w", err)
	}

	err = validator.Validator.Struct(cfg)
	if err != nil {
		return cfg, fmt.Errorf("error during validate config: %w", err)
	}

	return cfg, nil
}
