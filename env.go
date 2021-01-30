package common

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	dir, err := os.Getwd()
	path := dir + "/.env"

	if err != nil {
		log.Println("Not able to get current working director")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir = filepath.Dir(dir)
		path = dir + "/.env"
	}

	// loads values from .env into the system
	log.Println("Loading variables from " + path)

	if err := godotenv.Load(path); err != nil {
		log.Println("No .env file found in path " + path)
	}
}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

// LoadEnvVars will load a ".env[.development|.test]" file if it exists and set ENV vars.
// Useful in development and test modes. Not used in production.
func LoadEnvVars() {
	env := GetEnv("GIN_ENV", "development")

	if env == "production" || env == "staging" {
		log.Println("Not using .env file in production or staging.")
		return
	}

	filename := ".env." + env

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = ".env"
	}

	err := godotenv.Load(filename)
	if err != nil {
		log.Println(".env file not loaded")
	}
}
