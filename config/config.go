package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpires  int
	Environment string
}

func Load() *Config {
	// Load .env nếu có (local dev)
	// Docker sẽ dùng environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse JWT expires
	jwtStr := getEnv("JWT_EXPIRES_HOURS", "24") // ← FIX typo: EXPiRES → EXPIRES
	jwtInt, err := strconv.Atoi(jwtStr)
	if err != nil {
		jwtInt = 24
	}

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: buildDatabaseURL(),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		JWTExpires:  jwtInt,
		Environment: getEnv("ENV", "development"),
	}

	// Log config (mask sensitive data)
	log.Printf("✅ Config loaded:")
	log.Printf("   Port: %s", cfg.Port)
	log.Printf("   Environment: %s", cfg.Environment)
	log.Printf("   JWT Expires: %d hours", cfg.JWTExpires)
	log.Printf("   Database: %s", maskDatabaseURL(cfg.DatabaseURL))

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDatabaseURL() string {
	// ✅ Kiểm tra nếu có DATABASE_URL đầy đủ (Docker/Production)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	// ✅ Nếu không, build từ individual vars (Local development)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "inventory_pro")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// ✅ PostgreSQL connection string format
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

func maskDatabaseURL(url string) string {
	// Simple masking for logging
	return "postgres://***:***@<host>/<dbname>"
}
