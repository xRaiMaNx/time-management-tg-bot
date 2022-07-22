package timeconv

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var ErrBadWeekDay = errors.New("incorrect content of the days of the week field")
var ErrBadTime = errors.New("incorrect content of the time field")

var dayNums = map[string]int{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
}

func MinuteToStr(minutes int) string {
	hours := minutes / 60
	minutes %= 60
	hoursStr := strconv.Itoa(hours)
	if len(hoursStr) == 1 {
		hoursStr = "0" + hoursStr
	}

	minutesStr := strconv.Itoa(minutes)
	if len(minutesStr) == 1 {
		minutesStr = "0" + minutesStr
	}

	return hoursStr + ":" + minutesStr
}

func GetWeekStat() (curWeekDay int, curMinutes int) {
	curWeekDay = int(time.Now().Weekday())
	curMinutes = time.Now().Hour()*60 + time.Now().Minute()
	return
}

func GetScheduleWithoutTime(weekDays []string) (res [7]bool, e error) {
	for _, weekDay := range weekDays {
		if strings.ToLower(weekDay) == "all" {
			for i := 0; i < 7; i++ {
				res[i] = true
			}
			return
		}
		if strings.ToLower(weekDay) == "today" {
			res[int(time.Now().Weekday())] = true
			continue
		}
		if dayNum, ok := dayNums[strings.ToLower(weekDay)]; ok {
			res[dayNum] = true
		} else {
			e = ErrBadWeekDay
			return
		}
	}
	return
}

func GetScheduleEvent(weekDays []string, strTime string) (res [7]int, e error) {
	if len(weekDays) == 0 {
		e = ErrBadWeekDay
		return
	}
	minutes, e := getMinutesFromString(strTime)
	if e != nil {
		return
	}

	resBool, e := GetScheduleWithoutTime(weekDays)
	if e != nil {
		return
	}

	for i, value := range resBool {
		if value {
			res[i] = minutes
		} else {
			res[i] = -1
		}
	}
	return
}

func getMinutesFromString(strTime string) (int, error) {
	times := strings.Split(strTime, ":")
	if len(times) != 2 {
		return 0, ErrBadTime
	}

	if len(times[1]) != 2 {
		return 0, ErrBadTime
	}
	minutes, err := strconv.Atoi(times[1])
	if err != nil || minutes < 0 || minutes > 60 {
		return 0, ErrBadTime
	}

	if len(times[0]) > 2 {
		return 0, ErrBadTime
	}
	hours, err := strconv.Atoi(times[0])
	if err != nil || hours < 0 || hours > 24 {
		return 0, ErrBadTime
	}
	return hours*60 + minutes, nil
}
