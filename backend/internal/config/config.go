package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Port int
	Env  string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") == "development" {
			log.Printf("! .env not loaded: %v", err)
		}
	}

	pgPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}

	appPort, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	return &Config{
			App: AppConfig{
				Port: appPort,
				Env:  getEnv("APP_ENV", "development"),
			},
			Postgres: PostgresConfig{
				Host:     getEnv("POSTGRES_HOST", "localhost"),
				Port:     pgPort,
				User:     getEnv("POSTGRES_USER", "postgres"),
				Password: getEnv("POSTGRES_PASSWORD", "postgres"),
				DB:       getEnv("POSTGRES_DB", "workspace"),
			},
			Redis: RedisConfig{
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     redisPort,
				Password: getEnv("REDIS_PASSWORD", ""),
			},
		},
		nil
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}
	return val
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		url.QueryEscape(c.Password),
		c.DB,
		getEnv("POSTGRES_SSLMODE", "disable"),
	)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
