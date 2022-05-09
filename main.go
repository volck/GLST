package main

import "os"

func main() {
	server := NewServer()
	server.Routes()
	server.steamApiKey = os.Getenv("STEAM_WEBAPI")
	server.Run(":1337")
}