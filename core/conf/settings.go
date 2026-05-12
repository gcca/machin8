package conf

import (
	"errors"
	"os"
)

type Config struct {
	Port   string
	Secret string
	DBUrl  string
}

var (
	Settings Config
)

func InitSettings() error {
	Settings = Config{
		Port:   env("PORT", "8000"),
		Secret: env("SECRET", ""),
		DBUrl:  env("DATABASE_URL", ""),
	}

	var errs []error
	if Settings.Secret == "" {
		errs = append(errs, errors.New("SECRET is required"))
	}
	if Settings.DBUrl == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}

	return errors.Join(errs...)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
