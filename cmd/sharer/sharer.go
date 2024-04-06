package main

import (
	"log"
	"os"

	"github.com/cloyop/sharer/client"
	"github.com/cloyop/sharer/server"
)

func main() {

	args := os.Args
	if len(args) < 2 {
		//Print Guide
		log.Fatal("not Enough arguments")
	}
	switch args[1] {
	case "server":
		server.RecieveHandler(args[2:])
		return
	case "client":
		client.SharerHandler(args[2:])
		return
	case "config":
		// ToDo
		return
	default:
		// PrintGuide
		log.Fatalf("invalid mode %v", args[1])
	}

}
