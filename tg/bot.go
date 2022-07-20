package tg

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"io/ioutil"
)

const (
	tgApi = "https://api.telegram.org/bot"
)

func Run(token string) {
	botUrl := tgApi + token
	for {
		updates, err := getUpdates(botUrl)
		if err != nil {
			log.Println("error in updates: ", err.Error())
		}
		fmt.Println(updates)
	}
}

func getUpdates(botUrl string) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates")
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