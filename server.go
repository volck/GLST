package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type glst struct {
	Steamid     string `db:"STEAMID" json:"steamid"`
	Appid       int    `db:"appid" json:"appid"`
	LoginToken  string `db:"login_token"  json:"login_token"`
	Memo        string `db:"memo" json:"memo"`
	IsDeleted   bool   `db:"is_deleted" json:"is_deleted"`
	IsExpired   bool   `db:"is_expired"  json:"is_expired"`
	IsUsed      bool   `db:"Is_used" json:"is_used"`
	RtLastLogon int    `db:"rt_last_logon" json:"rt_last_logon"`
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

	const file string = "glst.db"
	mydb, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("error opening database", err)
	}

	createdbifnoexists, err := mydb.Prepare("CREATE TABLE IF NOT EXISTS [glst](Steamid text PRIMARY KEY, AppId int, LoginToken text, Memo text, isDeleted int, isExpired int, isUsed int, rtLastLogon datetime );")
	if err != nil {
		fmt.Println("create table failed.")
	}

	res, err := createdbifnoexists.Exec()
	if err != nil {
		fmt.Println(err)
	}
	rowsaffected, err := res.RowsAffected()

	if rowsaffected != 0 {
		fmt.Printf("created database: rows affected: %v \n ", rowsaffected)
	}

	return &server{
		router: router,
		db:     mydb,
	}
}
