package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		User     string
		Password string
		Host     string
		Port     int
		Name     string
		DSN      string
	}
}

func LoadConfigFromEnv() (config Config, err error) {
	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get server port from environment variable
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		return config, fmt.Errorf("SERVER_PORT environment variable is not set")
	}
	config.Server.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return config, fmt.Errorf("invalid SERVER_PORT value: %v", err)
	}

	// Get database details from environment variables
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Name = os.Getenv("DB_NAME")

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		return config, fmt.Errorf("DB_PORT environment variable is not set")
	}
	config.Database.Port, err = strconv.Atoi(dbPortStr)
	if err != nil {
		return config, fmt.Errorf("invalid DB_PORT value: %v", err)
	}

	// Construct the DSN
	config.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	return config, nil
}
