package main

import (
	"log"
	"os"

	"github.com/mrbttf/bridge-server/pkg/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(server.Start(port))
}
