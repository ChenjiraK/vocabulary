package database

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	_ = godotenv.Load()

	return gorm.Open(postgres.Open(PostgresDSN()), &gorm.Config{})
}

func PostgresDSN() string {
	_ = godotenv.Load()

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "vocabulary"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)
}

func PostgresMigrationURL() string {
	_ = godotenv.Load()

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(getEnv("DB_USER", "postgres"), getEnv("DB_PASSWORD", "")),
		Host:   fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432")),
		Path:   getEnv("DB_NAME", "vocabulary"),
	}

	query := databaseURL.Query()
	query.Set("sslmode", getEnv("DB_SSLMODE", "disable"))
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
