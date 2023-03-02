package main

import (
	"os"

	"github.com/mrbttf/bridge-server/pkg/config"
	"github.com/mrbttf/bridge-server/pkg/core/services/session"
	"github.com/mrbttf/bridge-server/pkg/core/services/auth"
	"github.com/mrbttf/bridge-server/pkg/db"
	"github.com/mrbttf/bridge-server/pkg/log"
	"github.com/mrbttf/bridge-server/pkg/repositories"
	"github.com/mrbttf/bridge-server/pkg/server"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
		log.Info("Using default port 8080")
	}
	config, err := config.GetConfig(os.Getenv("ENV"))
	if err != nil {
		log.Fatal(err)
	}
	postgresDB, err := db.New(&config)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	repository := repositories.NewSessionRepository(postgresDB)
	playerRepository := repositories.NewPlayerRepository(postgresDB)
	serviceSession := session.New(
		repository,
		playerRepository,
	)
	userRepository := repositories.NewUserRepository(postgresDB)
	authService:= auth.New(
		userRepository,
	)
	server := server.New(serviceSession, authService, config)
	err = server.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
