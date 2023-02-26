package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerSubroute    string
}

func GetConfig(env string) (Config, error) {
	if env == "dev" {
		err := godotenv.Load()
		if err != nil {
			return Config{}, fmt.Errorf("Couldn't load config: %w", err)
		}
	}

	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerSubroute:     os.Getenv("SERVER_SUBROUTE"),
	}, nil
}
