
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	MinioHost      string
	MinioRootUser  string
	MinioRootPass  string
	RabbitMQUri    string
}

// New returns a new Config struct
func New() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:         getEnv("DB_HOST", ""),
		DBPort:         getEnv("DB_PORT", ""),
		DBUser:         getEnv("POSTGRES_USER", ""),
		DBPassword:     getEnv("POSTGRES_PASSWORD", ""),
		DBName:         getEnv("POSTGRES_DB", ""),
		MinioHost:      getEnv("MINIO_HOST", ""),
		MinioRootUser:  getEnv("MINIO_ROOT_USER", ""),
		MinioRootPass: getEnv("MINIO_ROOT_PASSWORD", ""),
		RabbitMQUri:    getEnv("RABBITMQ_URI", ""),
	}
}

// getEnv gets an environment variable or sets a default if not present
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
