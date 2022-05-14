package main

import (
	"fmt"
	"os"
)

func main() {
	server := NewServer()
	server.Routes()
	server.steamApiKey = os.Getenv("STEAM_WEBAPI")
	if server.steamApiKey != "" {
		server.gslList = server.getAllGsl(server.steamApiKey)
		go server.refreshExpiredTokens()
		go server.getAllGSLReccuring()
		server.Run(":1337")
	} else {
		fmt.Println("STEAM_WEBAPI not defined")
		os.Exit(1)
	}
}
