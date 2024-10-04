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
		Type  string
		MySQL struct {
			User     string
			Password string
			Host     string
			Port     int
			Name     string
			DSN      string
		}
		MongoDB struct {
			URI        string
			Database   string
			Collection string
		}
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

	// Get database type
	config.Database.Type = os.Getenv("DB_TYPE")
	if config.Database.Type == "" {
		return config, fmt.Errorf("DB_TYPE environment variable is not set")
	}

	switch config.Database.Type {
	case "mysql":
		err = loadMySQLConfig(&config)
	case "mongodb":
		err = loadMongoDBConfig(&config)
	default:
		return config, fmt.Errorf("unsupported DB_TYPE: %s", config.Database.Type)
	}

	return config, err
}

func loadMySQLConfig(config *Config) error {
	config.Database.MySQL.User = os.Getenv("MYSQL_USER")
	config.Database.MySQL.Password = os.Getenv("MYSQL_PASSWORD")
	config.Database.MySQL.Host = os.Getenv("MYSQL_HOST")
	config.Database.MySQL.Name = os.Getenv("MYSQL_NAME")

	dbPortStr := os.Getenv("MYSQL_PORT")
	if dbPortStr == "" {
		return fmt.Errorf("MYSQL_PORT environment variable is not set")
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return fmt.Errorf("invalid MYSQL_PORT value: %v", err)
	}
	config.Database.MySQL.Port = dbPort

	// Construct the DSN
	config.Database.MySQL.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Database.MySQL.User,
		config.Database.MySQL.Password,
		config.Database.MySQL.Host,
		config.Database.MySQL.Port,
		config.Database.MySQL.Name,
	)

	return nil
}

func loadMongoDBConfig(config *Config) error {
	config.Database.MongoDB.URI = os.Getenv("MONGODB_URI")
	if config.Database.MongoDB.URI == "" {
		return fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	config.Database.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	if config.Database.MongoDB.Database == "" {
		return fmt.Errorf("MONGODB_DATABASE environment variable is not set")
	}

	config.Database.MongoDB.Collection = os.Getenv("MONGODB_COLLECTION")
	if config.Database.MongoDB.Collection == "" {
		return fmt.Errorf("MONGODB_COLLECTION environment variable is not set")
	}

	return nil
}
