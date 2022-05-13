package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

type glst struct {
	Steamid     string `json:"steamid"`
	Appid       int    `json:"appid"`
	LoginToken  string `json:"login_token"`
	Memo        string `json:"memo"`
	IsDeleted   bool   `json:"is_deleted"`
	IsExpired   bool   `json:"is_expired"`
	IsUsed      bool   `json:"is_used"`
	RtLastLogon int    `json:"rt_last_logon"`
}

type steamServer struct {
	Response struct {
		Servers []glst `json:"servers"`
	} `json:"response"`
}

type server struct {
	router      *gin.Engine
	db          *sql.DB
	steamApiKey string
	gslList     steamServer
}

func (s *server) Run(port string) {
	log.Fatal(s.router.Run(port))
}

func NewServer() *server {

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.New()

	router.SetTrustedProxies([]string{""})

	router.Use(JSONLogMiddleware())
	router.Use(gin.Recovery())

	return &server{
		router: router,
	}

}
