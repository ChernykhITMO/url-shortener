package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env        string     `yaml:"env"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
	Service    Service    `yaml:"service"`
}

type Service struct {
	MaxAttempts int `yaml:"max_attempts"`
}

type HTTPServer struct {
	Address           string        `yaml:"address"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
}

type Postgres struct {
	DSN             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

func InitConfig(path string) (*Config, error) {
	const op = "config.InitConfig"
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	normalizationCfg(&cfg)

	return &cfg, nil
}

func normalizationCfg(cfg *Config) {
	if cfg.Env == "" {
		cfg.Env = "local"
	}

	if cfg.Service.MaxAttempts <= 0 {
		cfg.Service.MaxAttempts = 20
	}

	if cfg.HTTPServer.Address == "" {
		cfg.HTTPServer.Address = ":8080"
	}
}
