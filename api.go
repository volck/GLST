package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func (s *server) getAllGsl(steamWebapi string) steamServer {
	logMsg := fmt.Sprintf(`GLST [%s] GET all tokens`, time.Now().String())
	log.Info(logMsg)
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

func (s *server) getRandomToken(gslList steamServer) glst {
	chosenToken := glst{}
	for {
		randNr := rand.Intn(len(gslList.Response.Servers)-0) + 0
		token := gslList.Response.Servers[randNr]
		if !token.IsUsed && !token.IsExpired {
			chosenToken = token
			break
		}
	}
	return chosenToken

}

func (s *server) renewToken(theGLST glst) {

	body := strings.NewReader(fmt.Sprintf(`key=%s&steamid=%s`, s.steamApiKey, theGLST.Steamid))
	req, err := http.NewRequest("POST", "https://api.steampowered.com/IGameServersService/ResetLoginToken/v1/", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)

	}
	_, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	//_ := string(response)
	logMsg := fmt.Sprintf(`GLST [%s] UPDATE token %s`, time.Now().String(), theGLST.LoginToken)
	log.Info(logMsg)
}

func (s *server) refreshExpiredTokens() {
	for {
		for _, theGLST := range s.gslList.Response.Servers {
			if theGLST.IsExpired && !theGLST.IsUsed {
				s.renewToken(theGLST)
			}
		}
		sleepPeriod := 30 * time.Minute
		logMsg := fmt.Sprintf(`GLST [%s ] DONE updating token. Sleeping for %v`, time.Now().String(), sleepPeriod)
		log.Info(logMsg)
		s.gslList = s.getAllGsl(s.steamApiKey)
		time.Sleep(sleepPeriod)
	}
}
