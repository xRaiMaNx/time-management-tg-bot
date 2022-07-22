package tg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gammazero/deque"
	"github.com/xraimanx/time-management-tg-bot/timeconv"
)

type Schedule [7]map[int](map[int]string)

const (
	tgAPI       = "https://api.telegram.org/bot"
	threadCount = 8
)

func Run(token string) {
	var sch Schedule
	var sm sync.Mutex // schedule mutex
	for i := 0; i < 7; i++ {
		sch[i] = make(map[int](map[int]string))
	}

	botURL := tgAPI + token
	offset := 0

	go timeManage(botURL, &sch)

	var d deque.Deque[Update]
	var dm sync.Mutex // deque Mutex
	updateCh := make(chan struct{}, 128)
	for i := 0; i < threadCount-1; i++ {
		go updateProcessing(botURL, &d, &dm, updateCh, &sm, &sch)
	}

	log.Println("ok")

	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			log.Println("error in updates: ", err.Error())
		}
		for _, update := range updates {
			log.Println(update)
			dm.Lock()
			d.PushBack(update)
			dm.Unlock()
			updateCh <- struct{}{}
			offset = update.UpdateID + 1
		}
	}
}

func getUpdates(botURL string, offset int) ([]Update, error) {
	var resp *http.Response
	var body []byte
	var rr RestResponse
	var err error
	if resp, err = http.Get(botURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &rr); err != nil {
		return nil, err
	}
	return rr.Result, nil
}

func timeManage(botURL string, sch *Schedule) {
	for {
		<-timeconv.Wait()

		curWeekDay, curMinutes := timeconv.GetWeekStat()
		if _, ok := sch[curWeekDay][curMinutes]; !ok {
			continue
		}

		timeStr := timeconv.MinuteToStr(curMinutes)

		for chatID, event := range sch[curWeekDay][curMinutes] {
			var msg BotMsg

			msg.ChatID = chatID
			msg.Text = timeStr + " " + event

			if err := sendMessage(botURL, msg); err != nil {
				log.Println("error in time manage: ", err.Error())
			}
		}
	}
}

func updateProcessing(botURL string, d *deque.Deque[Update],
	dm *sync.Mutex, updateCh <-chan struct{}, sm *sync.Mutex, sch *Schedule) {
	for {
		<-updateCh
		dm.Lock()
		update := d.PopFront()
		log.Println("I got new update:", update)
		dm.Unlock()
		err := respond(botURL, update, sm, sch)
		if err != nil {
			log.Println("error in respond: ", err.Error())
		}
	}
}

func respond(botURL string, update Update, sm *sync.Mutex, sch *Schedule) error {
	chatID := update.Message.Chat.ChatID

	userMessage := update.Message.Text
	userMessage = CollapseSpaces(userMessage)
	userMessage = strings.Trim(userMessage, " ")
	userWords := strings.Split(userMessage, " ")
	log.Println("user words:", len(userWords), userWords)
	if len(userWords) == 0 {
		return errors.New("empty message")
	}

	var msg BotMsg
	switch userWords[0] {
	case "/start":
		log.Println("I got /start")
		msg = handleStart(chatID)
	case "/add_event":
		msg = handleNewEvent(chatID, sm, userWords[1:], sch)
	case "/delete_event":
		msg = handleDeleteEvent(chatID, sm, userWords[1:], sch)
	case "/show_schedule":
		msg = handleShowSchedule(chatID, userWords[1:], sch)
	default:
		return nil
	}
	err := sendMessage(botURL, msg)
	return err
}

func sendMessage(botURL string, msg BotMsg) error {
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
