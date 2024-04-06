package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cloyop/sharer/client"
	"github.com/cloyop/sharer/server"
)

func main() {

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not Enought Arguments")
		fmt.Println(client.ArgumentsGuide() + server.ArgumentsGuide())
		return
	}
	switch args[1] {
	case "server":
		server.RecieveHandler(server.ParseArguments(args[2:]))
		return
	case "client":
		client.SharerHandler(client.ParseArguments(args[2:]))
		return
	default:
		log.Fatalf("invalid Argument %v", args[1])
		fmt.Println(client.ArgumentsGuide() + server.ArgumentsGuide())
	}

}
