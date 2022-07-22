package tg

import (
	"sync"
	"time"

	"github.com/xraimanx/time-management-tg-bot/timeconv"
)

func handleStart(chatID int) BotMsg {
	text := `/show_schedule [days] - Check your schedule
/add_event [days] [time] [event] - Add new event
/delete_event [days] [time] [event] - Delete event
days - list of week days("Monday", "Tuesday", etc. or "all", "today")
time - time of the event(for example, 06:31 or 8:22)
event - name of the event
Examples: 
add_event Monday Friday 21:00 Relaxed time
add_event all 07:30 Breakfast`
	return BotMsg{ChatID: chatID, Text: text}
}

func handleNewEvent(chatID int, sm *sync.Mutex,
	userWords []string, sch *Schedule) (msg BotMsg) {
	msg.ChatID = chatID
	n := len(userWords)
	if n < 3 {
		msg.Text = "Not enough arguments for command /add_event"
		return
	}

	if len(userWords[n-1]) > 40 {
		msg.Text = "The event name is either too big or empty"
		return
	}

	days, err := timeconv.GetScheduleEvent(userWords[:n-2], userWords[n-2])
	if err != nil {
		msg.Text = err.Error()
		return
	}

	var isWasSmth bool
	sm.Lock()
	for i := 0; i < 7; i++ {
		if days[i] != -1 {
			if sch[i][days[i]] == nil {
				sch[i][days[i]] = make(map[int]string)
			}
			if _, ok := sch[i][days[i]][chatID]; ok {
				isWasSmth = true
			}
			sch[i][days[i]][chatID] = userWords[n-1]
		}
	}
	sm.Unlock()

	if isWasSmth {
		msg.Text = "Done, but I changed old event/events, which should be in the same time"
	} else {
		msg.Text = "Done"
	}
	return
}

func handleDeleteEvent(chatID int, sm *sync.Mutex,
	userWords []string, sch *Schedule) (msg BotMsg) {
	msg.ChatID = chatID
	n := len(userWords)
	if n < 3 {
		msg.Text = "Not enough arguments for command /delete_event"
		return
	}

	if len(userWords[n-1]) > 40 {
		msg.Text = "The event name is either too big or empty"
		return
	}

	days, err := timeconv.GetScheduleEvent(userWords[:n-2], userWords[n-2])
	if err != nil {
		msg.Text = err.Error()
		return
	}

	var isDeletedSmth bool
	sm.Lock()
	for i := 0; i < 7; i++ {
		if days[i] != -1 {
			if sch[i][days[i]] == nil {
				continue
			}
			if sch[i][days[i]][chatID] == userWords[n-1] {
				delete(sch[i][days[i]], chatID)
				isDeletedSmth = true
			}
		}
	}
	sm.Unlock()
	if !isDeletedSmth {
		msg.Text = "Nothing to delete"
	} else {
		msg.Text = "Done"
	}
	return
}

func handleShowSchedule(chatID int, userWords []string, sch *Schedule) (msg BotMsg) {
	msg.ChatID = chatID
	weekDaysBool, err := timeconv.GetScheduleWithoutTime(userWords)
	if err != nil {
		msg.Text = err.Error()
		return
	}

	for dayNum, isNeed := range weekDaysBool {
		if !isNeed {
			continue
		}
		msg.Text += time.Weekday(dayNum).String() + ":\n"
		for i := 0; i < 24*60; i++ {
			if _, ok := sch[dayNum][i]; !ok {
				continue
			}
			event, ok := sch[dayNum][i][chatID]
			if !ok {
				continue
			}
			timeStr := timeconv.MinuteToStr(i)
			msg.Text += timeStr + " " + event + "\n"
		}
		msg.Text += "\n"
	}
	return
}
