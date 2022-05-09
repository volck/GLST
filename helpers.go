package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func SteamPoweredAPIRequest(endpoint string, webapikey string) (jsonString string) {
	key := "key=" + webapikey
	url := "https://api.steampowered.com/" + endpoint + "?&format=json&" + key

	Client := http.Client{
		Timeout: time.Second * 2, //Maximum of 2 secs
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
