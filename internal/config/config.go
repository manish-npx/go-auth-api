package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string   `yaml:"env"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Security Security `yaml:"security"`
	Logging  Logging  `yaml:"logging"`
}

type Server struct {
	Port                int `yaml:"port"`
	ReadTimeoutSeconds  int `yaml:"readTimeoutSeconds"`
	WriteTimeoutSeconds int `yaml:"writeTimeoutSeconds"`
}

type Database struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	User                   string `yaml:"user"`
	Pass                   string `yaml:"pass"`
	Name                   string `yaml:"name"`
	MaxConns               int32  `yaml:"maxConns"`
	MinConns               int32  `yaml:"minConns"`
	MaxConnLifetimeMinutes int    `yaml:"maxConnLifetimeMinutes"`
}

type Security struct {
	JWTSecret string `yaml:"jwtSecret"`
}

type Logging struct {
	Level string `yaml:"level"`
	JSON  bool   `yaml:"json"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	if cfg.Security.JWTSecret == "" {
		return nil, fmt.Errorf("security.jwtSecret is required")
	}
	return &cfg, nil
}

func ConfigureLogger(cfg *Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	if !cfg.Logging.JSON {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	}
	switch cfg.Logging.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	return log.Logger
}
