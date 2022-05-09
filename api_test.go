package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func Test_server_UnExpireGLSTinDatabase(t *testing.T) {

	const file string = "glst.db"
	mydb, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("error opening database", err)
	}

	s := server{
		steamApiKey: os.Getenv("STEAM_WEBAPI"),
		db:          mydb,
	}

	allGLSTS := getAllGsl(s.steamApiKey)

	for _, glst := range allGLSTS.Response.Servers {
		if glst.IsExpired {
			fmt.Printf("%s is expired, trying this\n", glst.Steamid)
			s.renewToken(glst.Steamid)
			s.UnExpireGLSTinDatabase(glst.Steamid)
		}
	}

}

func Test_server_INSERTglsts(t *testing.T) {

	const file string = "glst.db"
	mydb, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("error opening database", err)
	}

	s := server{
		steamApiKey: os.Getenv("STEAM_WEBAPI"),
		db:          mydb,
	}

	s.INSERTglsts()

}

func Test_getNonExpiredAndNonUsedToken(t *testing.T) {

	const file string = "glst.db"
	mydb, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("error opening database", err)
	}

	s := server{
		steamApiKey: os.Getenv("STEAM_WEBAPI"),
		db:          mydb,
	}
	GetFreshToken(s.db)

}

func Test_setTokenToUnused(t *testing.T) {

	const file string = "glst.db"
	mydb, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("error opening database", err)
	}

	s := server{
		steamApiKey: os.Getenv("STEAM_WEBAPI"),
		db:          mydb,
	}

	allGLSTS := getAllGsl(s.steamApiKey)

	setTokenToUnused(s.db, allGLSTS.Response.Servers[0].Steamid)
}

func Test_setTokenToUsed(t *testing.T) {

	const file string = "glst.db"
	mydb, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("error opening database", err)
	}

	s := server{
		steamApiKey: os.Getenv("STEAM_WEBAPI"),
		db:          mydb,
	}

	allGLSTS := getAllGsl(s.steamApiKey)

	setTokenToUsed(s.db, allGLSTS.Response.Servers[0].Steamid)

}
