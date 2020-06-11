package main

import (
	"log"

	"github.com/mrbttf/bridge-server/pkg/server"
)

func main() {
	log.Fatal(server.Start("80"))
}
