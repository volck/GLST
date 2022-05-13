package main

import "os"

func main() {
	server := NewServer()
	server.Routes()
	server.steamApiKey = os.Getenv("STEAM_WEBAPI")
	server.gslList = server.getAllGsl(server.steamApiKey)
	go server.refreshExpiredTokens()
	server.Run(":1337")
}