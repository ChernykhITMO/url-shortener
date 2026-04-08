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
	AliasLength int `yaml:"alias_length"`
}

type HTTPServer struct {
	Address           string        `yaml:"address"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
	MaxBodyBytes      int64         `yaml:"max_body_bytes"`
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

	normalizeConfig(&cfg)
	applyEnv(&cfg)

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cfg, nil
}

func normalizeConfig(cfg *Config) {
	if cfg.Env == "" {
		cfg.Env = "local"
	}

	if cfg.Service.MaxAttempts <= 0 {
		cfg.Service.MaxAttempts = 20
	}

	if cfg.Service.AliasLength <= 0 {
		cfg.Service.AliasLength = 10
	}

	if cfg.HTTPServer.Address == "" {
		cfg.HTTPServer.Address = ":8080"
	}

	if cfg.HTTPServer.MaxBodyBytes <= 0 {
		cfg.HTTPServer.MaxBodyBytes = 1024
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		cfg.HTTPServer.Address = v
	}

	if v := os.Getenv("POSTGRES_DSN"); v != "" {
		cfg.Postgres.DSN = v
	}
}

func validateConfig(cfg *Config) error {
	const op = "config.validateConfig"

	if cfg.HTTPServer.Address == "" {
		return fmt.Errorf("%s: http_server.address is required", op)
	}

	if cfg.Service.MaxAttempts <= 0 {
		return fmt.Errorf("%s: service.max_attempts must be > 0", op)
	}

	if cfg.Service.AliasLength <= 0 {
		return fmt.Errorf("%s: service.alias_length must be > 0", op)
	}

	if err := validateHTTPServerConfig(cfg.HTTPServer); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := validatePostgresConfig(cfg.Postgres); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func validateHTTPServerConfig(c HTTPServer) error {
	const op = "config.validateHTTPServerConfig"

	if c.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("%s: http_server.read_header_timeout must be > 0", op)
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("%s: http_server.read_timeout must be > 0", op)
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("%s: http_server.write_timeout must be > 0", op)
	}

	if c.IdleTimeout <= 0 {
		return fmt.Errorf("%s: http_server.idle_timeout must be > 0", op)
	}

	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("%s: http_server.shutdown_timeout must be > 0", op)
	}

	if c.MaxBodyBytes <= 0 {
		return fmt.Errorf("%s: http_server.max_body_bytes must be > 0", op)
	}

	return nil
}

func validatePostgresConfig(c Postgres) error {
	const op = "config.validatePostgresConfig"

	if c.DSN != "" {
		if c.MaxOpenConns < 0 {
			return fmt.Errorf("%s: postgres.max_open_conns must be >= 0", op)
		}

		if c.MaxIdleConns < 0 {
			return fmt.Errorf("%s: postgres.max_idle_conns must be >= 0", op)
		}

		if c.MaxOpenConns > 0 && c.MaxIdleConns > c.MaxOpenConns {
			return fmt.Errorf("%s: postgres.max_idle_conns must be <= postgres.max_open_conns", op)
		}

		if c.ConnMaxLifetime < 0 {
			return fmt.Errorf("%s: postgres.conn_max_lifetime must be >= 0", op)
		}
	}

	return nil
}
