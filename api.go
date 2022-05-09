package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"strings"
)

func GETglstsFromDB(databaseConnection *sql.DB) (steamServer, error) {
	row := databaseConnection.QueryRow("select * from glst")
	fmt.Println(row)
	myResponseSteamServer := steamServer{}

	return myResponseSteamServer, nil
}

func GetFreshToken(databaseConnection *sql.DB) (glst, error) {
	rows := databaseConnection.QueryRow("select * from glst where isexpired != 1 AND isused != 1")
	myResponseSteamServer := glst{}
	err := rows.Scan(&myResponseSteamServer.Steamid, &myResponseSteamServer.Appid, &myResponseSteamServer.LoginToken,
		&myResponseSteamServer.Memo, &myResponseSteamServer.IsDeleted,
		&myResponseSteamServer.IsExpired, &myResponseSteamServer.IsUsed, &myResponseSteamServer.RtLastLogon)
	if err != nil {
		fmt.Printf("scan err: %s\n", err.Error())
	}
	if err.Error() == "sql: no rows in result set" {
		return myResponseSteamServer, err
	} else {
		setTokenToUsed(databaseConnection, myResponseSteamServer.LoginToken)
		return myResponseSteamServer, nil
	}
}

func setTokenToUsed(databaseConnection *sql.DB, steamID string) {
	result, err := databaseConnection.ExecContext(context.Background(), `UPDATE glst SET isUsed = true WHERE STEAMID = ?;`, steamID)
	if err != nil {
		fmt.Printf("execcontext failed %v\n", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("rows affected errored: %v\n", err)
	}
	if rows != 1 {
		fmt.Printf("expected single row affected, got %d rows affected\n", rows)
	}
	fmt.Printf("setTokenToUsed: %v rows affected\n", rows)
}

func setTokenToUnused(databaseConnection *sql.DB, steamID string) {

	result, err := databaseConnection.ExecContext(context.Background(), `UPDATE glst SET isUsed = false WHERE steamid = ?;`, steamID)
	if err != nil {
		fmt.Printf("execcontext failed %v\n", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("rows affected errored: %v\n", err)
	}
	if rows != 1 {
		fmt.Printf("expected single row affected, got %d rows affected\n", rows)
	}
	fmt.Printf("setTokenToUsed: %v rows affected\n", rows)
}

func (s *server) INSERTglsts() {
	jsonString := SteamPoweredAPIRequest("IGameServersService/GetAccountList/v1/", s.steamApiKey)

	var serverentry steamServer
	err := json.Unmarshal([]byte(jsonString), &serverentry)
	if err != nil {
		fmt.Println("error-> ", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		fmt.Println("db error", err)
	}

	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	allKeys := getAllGsl(s.steamApiKey)

	stmt, err := tx.Prepare("INSERT INTO glst(Steamid, Appid, LoginToken, Memo, isExpired, isUsed, rtLastLogon) VALUES( ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	for _, key := range allKeys.Response.Servers {

		if _, err := stmt.Exec(key.Steamid, key.Appid, key.LoginToken, key.Memo, key.IsExpired, key.IsUsed, key.RtLastLogon); err != nil {
			fmt.Printf("exec err: %v\n", err)
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Printf("commit err: %v\n", err)
	}

}

func (s *server) UnExpireGLSTinDatabase(steamid string) {

	tx, err := s.db.Begin()
	if err != nil {
		fmt.Println("db error", err)
	}

	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
	stmt, err := tx.Prepare("UPDATE glst SET isexpired = 0, isdeleted = 0 WHERE STEAMID = ? ")
	if err != nil {
		fmt.Println(err)
	}

	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	if _, err := stmt.Exec(steamid); err != nil {
		fmt.Println("exec err", err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println("commit err", err)
	}

}

func (s *server) renewToken(token string) (string, error) {

	body := strings.NewReader(fmt.Sprintf(`key=%s&steamid=%s`, s.steamApiKey, token))
	req, err := http.NewRequest("POST", "https://api.steampowered.com/IGameServersService/ResetLoginToken/v1/", body)
	if err != nil {
		fmt.Printf("newrequest failed: %v", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("http.defaultclient.do failed: %v\n", err)
	}
	response, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("readAll: %v\n", readErr)
	}
	responseStr := string(response)
	fmt.Printf("responsestring: %v\n", responseStr)

	return token, err

}

func getAllGsl(steamWebapi string) steamServer {

	jsonString := SteamPoweredAPIRequest("IGameServersService/GetAccountList/v1/", steamWebapi)

	var serverentry steamServer
	err := json.Unmarshal([]byte(jsonString), &serverentry)
	if err != nil {
		fmt.Printf("getAllGsl error: %v\n", err)
	} else {
		return serverentry
	}

	return serverentry
}

func (s *server) getUsedTokensFromDb() (glst, error) {
	rows := s.db.QueryRow("select * from glst where isused != 1")
	myResponseSteamServer := glst{}
	err := rows.Scan(&myResponseSteamServer.Steamid, &myResponseSteamServer.Appid, &myResponseSteamServer.LoginToken,
		&myResponseSteamServer.Memo, &myResponseSteamServer.IsDeleted,
		&myResponseSteamServer.IsExpired, &myResponseSteamServer.IsUsed, &myResponseSteamServer.RtLastLogon)
	if err != nil {
		fmt.Printf("scan err: %s\n", err.Error())
	}
	if err.Error() == "sql: no rows in result set" {
		return myResponseSteamServer, err
	} else {
		setTokenToUsed(s.db, myResponseSteamServer.LoginToken)
		return myResponseSteamServer, nil
	}

}
