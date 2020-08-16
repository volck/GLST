package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var flagvar string

type steamServer struct {
	Response struct {
		Servers []struct {
			Steamid     string `json:"steamid"`
			Appid       int    `json:"appid"`
			LoginToken  string `json:"login_token"`
			Memo        string `json:"memo"`
			IsDeleted   bool   `json:"is_deleted"`
			IsExpired   bool   `json:"is_expired"`
			IsUsed   bool      `json:"Is_used"`
			RtLastLogon int    `json:"rt_last_logon"`
		} `json:"servers"`
	} `json:"response"`
}

//awesome: https://mholt.github.io/json-to-go/
func doGetRequest(endpoint string, webapikey string)(jsonString string){
	key := "key=" + webapikey
 	url := "https://api.steampowered.com/" + endpoint + "?&format=json&" +  key
   
 	 Client := http.Client{
 	 	Timeout: time.Second * 2,  //Maximum of 2 secs
 	 }

 	 req, err := http.NewRequest(http.MethodGet, url, nil)
 	 if err != nil {
 	 	log.Fatal(err)
 	 }

 	 res, getErr := Client.Do(req)
 	 if getErr != nil {
 	 	log.Fatal(getErr)
 	 }

 	 body, readErr := ioutil.ReadAll(res.Body)
 	 if readErr != nil {
 	 	log.Fatal(readErr)
 	 }

 	 jsonString = string(body)

return 
} 

func renewToken(webapi string, token string){

	body := strings.NewReader(fmt.Sprintf(`key=%s&steamid=%s`,webapi,token))
	req, err := http.NewRequest("POST", "https://api.steampowered.com/IGameServersService/ResetLoginToken/v1/", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	response, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonString := string(response)
	fmt.Printf("[*] renewed: %s [*]\n", jsonString)
	return

}




func getAllGsl(webapikey string)(steamServer){
	

   
	// fmt.Print(string(body))
	jsonString := doGetRequest("IGameServersService/GetAccountList/v1/", webapikey)
	// fmt.Print(jsonString)
	 
	var serverentry steamServer
	err := json.Unmarshal([]byte(jsonString), &serverentry)
	if err != nil{
		fmt.Println("error-> ", err )
	} else {
		return serverentry
	} 

return serverentry
	}



func PrintAllExpiredGsls(webapikey string)(gsls steamServer){
	var list = getAllGsl(webapikey)
	var i = 0
	 for i <= len(list.Response.Servers)-1 {
	 	if list.Response.Servers[i].IsExpired{
	 	fmt.Printf("expired: %t \t steamid: %s \t login token: %s \t last_used: %d \n", list.Response.Servers[i].IsExpired, list.Response.Servers[i].Steamid, list.Response.Servers[i].LoginToken, list.Response.Servers[i].RtLastLogon)
	 	}
	 	i++;
		 }
	return list
}



func renewAllTokens(webapikey string)(newList steamServer){
	var list = getAllGsl(webapikey)
	var i = 0
	for i <= len(list.Response.Servers)-1 {
		if list.Response.Servers[i].IsExpired{
			renewToken(webapikey, list.Response.Servers[i].Steamid)
		}
		i++
	}
	newList = getAllGsl(webapikey)
	return newList



}



func main() {
    val, present := os.LookupEnv("steam_api")
		if !present {
			fmt.Println("[*] key not valid [*] ")
		} else {
		var used []string
		list := renewAllTokens(val)

		 http.HandleFunc("/NewToken", func (w http.ResponseWriter, r *http.Request) {
			 for {
			 choice := list.Response.Servers[rand.Intn(len(list.Response.Servers))]
			 if !choice.IsDeleted && !choice.IsExpired && !choice.IsUsed {
			 choice.IsUsed = true 
			 json.NewEncoder(w).Encode(choice.LoginToken)
			 used = append(used, choice.Steamid)
			 fmt.Printf("sent token for game: %s. Current list of tokens: %v\n",choice.Steamid,used)
			 break 
			 }
			}
		})


		fmt.Println("[*] listening on port 1337[*]")
		http.ListenAndServe(":1337", nil)
		 } 

	}

