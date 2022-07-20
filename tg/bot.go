package tg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	tgAPI = "https://api.telegram.org/bot"
)

func Run(token string) {
	botURL := tgAPI + token
	for {
		updates, err := getUpdates(botURL)
		if err != nil {
			log.Println("error in updates: ", err.Error())
		}
		fmt.Println(updates)
	}
}

func getUpdates(botURL string) ([]Update, error) {
	resp, err := http.Get(botURL + "/getUpdates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rr RestResponse
	err = json.Unmarshal(body, &rr)
	if err != nil {
		return nil, err
	}
	return rr.Result, nil
}
