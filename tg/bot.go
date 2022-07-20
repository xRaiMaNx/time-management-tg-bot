package tg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	tgAPI = "https://api.telegram.org/bot"
)

func Run(token string) {
	botURL := tgAPI + token
	offset := 0
	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			log.Println("error in updates: ", err.Error())
		}
		for _, update := range updates {
			err = respond(botURL, update)
			if err != nil {
				log.Println("error in response: ", err.Error())
			}
			offset = update.UpdateID + 1
		}
	}
}

func getUpdates(botURL string, offset int) ([]Update, error) {
	resp, err := http.Get(botURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
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

func respond(botURL string, update Update) error {
	var msg BotMessage
	msg.ChatID = update.Message.Chat.ChatID
	msg.Text = update.Message.Text
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
