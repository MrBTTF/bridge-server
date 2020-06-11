package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrbttf/bridge-server/pkg/server"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
		fmt.Println("using default port 8080")
	}
	log.Fatal(server.Start(port))
}
