package config

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	AppEnv              string        `env:"APP_ENV" envDefault:"local" validate:"oneof=local dev prod test"`
	MySQLDSN            string        `env:"MYSQL_DSN,required" validate:"required"`
	HTTPPort            string        `env:"HTTP_PORT" envDefault:"8080" validate:"required"`
	HTTPReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s" validate:"gt=0"`
	HTTPWriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s" validate:"gt=0"`
	HTTPIdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s" validate:"gt=0"`
	HTTPShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s" validate:"gt=0"`
	DBMaxOpenConns      int           `env:"DB_MAX_OPEN_CONNS" envDefault:"10" validate:"gt=0"`
	DBMaxIdleConns      int           `env:"DB_MAX_IDLE_CONNS" envDefault:"10" validate:"gte=0"`
	DBConnMaxLifetime   time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m" validate:"gt=0"`
	DBConnMaxIdleTime   time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m" validate:"gt=0"`
	LogLevelText        string        `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
}

func Load() (Config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse env: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c Config) HTTPAddr() string {
	if strings.HasPrefix(c.HTTPPort, ":") {
		return c.HTTPPort
	}

	return ":" + c.HTTPPort
}

func (c Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LogLevelText) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func loadDotEnv(path string) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open .env: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid .env line: %s", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			return fmt.Errorf("invalid .env key")
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan .env: %w", err)
	}

	return nil
}
