package timeconv

import (
	"errors"
	"strconv"
	"strings"
)

var ErrBadWeekDay = errors.New("incorrect content of the days of the week field")
var ErrBadTime = errors.New("incorrect content of the time field")

var dayNums = map[string]int{
	"monday":    0,
	"tuesday":   1,
	"wednesday": 2,
	"thursday":  3,
	"friday":    4,
	"saturday":  5,
	"sunday":    6,
}

func getMinutesFromString(time string) (int, error) {
	times := strings.Split(time, ":")
	if len(times) != 2 {
		return 0, ErrBadTime
	}

	minutes, err := strconv.Atoi(times[0])
	if err != nil {
		return 0, ErrBadTime
	}

	hours, err := strconv.Atoi(times[1])
	if err != nil {
		return 0, ErrBadTime
	}

	return hours*60 + minutes, nil
}

func GetSchedule(weekDays []string, time string) (res [7]int, e error) {
	if len(weekDays) == 0 {
		e = ErrBadWeekDay
		return
	}
	minutes, e := getMinutesFromString(time)
	if e != nil {
		return
	}
	for _, weedDay := range weekDays {
		if dayNum, ok := dayNums[weedDay]; ok {
			res[dayNum] = minutes
		} else {
			e = ErrBadWeekDay
			return
		}
	}
	return
}
