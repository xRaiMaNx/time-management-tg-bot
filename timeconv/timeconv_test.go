package timeconv

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchedule_Simple(t *testing.T) {
	for _, test := range []struct {
		name     string
		weekDays []string
		strTime  string
		want     [7]int
	}{
		{"monday", []string{"monday"}, "03:00", [7]int{-1, 3 * 60, -1, -1, -1, -1, -1}},
		{"tuesday", []string{"tuesday"}, "03:25", [7]int{-1, -1, 3*60 + 25, -1, -1, -1, -1}},
		{"wednesday", []string{"wednesday"}, "00:00", [7]int{-1, -1, -1, 0, -1, -1, -1}},
		{"thursday", []string{"thursday"}, "01:00", [7]int{-1, -1, -1, -1, 60, -1, -1}},
		{"friday", []string{"friday"}, "01:30", [7]int{-1, -1, -1, -1, -1, 90, -1}},
		{"saturday", []string{"saturday"}, "00:01", [7]int{-1, -1, -1, -1, -1, -1, 1}},
		{"sunday", []string{"sunday"}, "00:02", [7]int{2, -1, -1, -1, -1, -1, -1}},
		{"double_weekday", []string{"sunday", "sunday"}, "00:02", [7]int{2, -1, -1, -1, -1, -1, -1}},
		{"several_days", []string{"monday", "friday"}, "16:00", [7]int{-1, 16 * 60, -1, -1, -1, 16 * 60, -1}},
		{"several_days", []string{"monday", "tuesday", "friday", "sunday"}, "14:00", [7]int{14 * 60, 14 * 60, 14 * 60, -1, -1, 14 * 60, -1}},
		{"bad_order", []string{"friday", "monday"}, "11:33", [7]int{-1, 11*60 + 33, -1, -1, -1, 11*60 + 33, -1}},
		{"all_days", []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, "00:01", [7]int{1, 1, 1, 1, 1, 1, 1}},
		{"without_leading_zero_hour", []string{"monday"}, "3:00", [7]int{-1, 3 * 60, -1, -1, -1, -1, -1}},
		{"not_lower_case", []string{"Monday"}, "00:01", [7]int{-1, 1, -1, -1, -1, -1, -1}},
		{"all", []string{"all"}, "00:01", [7]int{1, 1, 1, 1, 1, 1, 1}},
		{"all_isnt_only_one", []string{"monday", "friday", "all"}, "00:01", [7]int{1, 1, 1, 1, 1, 1, 1}},
		{"one_zero", []string{"sunday"}, "00:00", [7]int{0, -1, -1, -1, -1, -1, -1}},
		{"all_zero", []string{"all"}, "00:00", [7]int{0, 0, 0, 0, 0, 0, 0}},
	} {
		t.Run(test.name, func(t *testing.T) {
			res, err := GetScheduleEvent(test.weekDays, test.strTime)
			require.NoError(t, err)

			assert.Equal(t, test.want, res)
		})
	}
}

func TestSchedule_IncorrectWeekDays(t *testing.T) {
	for _, test := range []struct {
		name     string
		weekDays []string
	}{
		{"empty", []string{}},
		{"monday", []string{"monda"}},
		{"num_instead_of_str", []string{"0"}},
		{"only_1_day", []string{"monday", "tuesday", "frida"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetScheduleEvent(test.weekDays, "00:01")
			assert.Error(t, err)

			if !errors.Is(err, ErrBadWeekDay) {
				t.Errorf("error isn't ErrBadWeekDay")
			}
		})
	}
}

func TestSchedule_IncorrectTime(t *testing.T) {
	for _, test := range []struct {
		name    string
		strTime string
	}{
		{"empty", ""},
		{"len_1", "01"},
		{"neg_hours", "-01:03"},
		{"neg_minutes", "01:-01"},
		{"minutes_more_than_60", "01:63"},
		{"hours_more_than_24", "25:01"},
		{"strange_symbols_after", "00:01k"},
		{"strange_symbols_before", "00:k01"},
		{"strange_symbols_after", "00ks:01"},
		{"strange_symbols_before", "k00:01"},
		{"seconds", "00:01:20"},
		{"without_leading_zero_minute", "00:1"},
		{"huge_hours", "9999999999999999999999999999999999:00"},
		{"huge_minutes", "00:9999999999999999999999999999999"},
		{"extra_zeros_hours", "000:01"},
		{"extra_zeros_minutes", "00:001"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetScheduleEvent([]string{"monday"}, test.strTime)
			assert.Error(t, err)

			if !errors.Is(err, ErrBadTime) {
				t.Errorf("error isn't ErrBadTime")
			}
		})
	}
}

func TestWait(t *testing.T) {
	ch := make(chan struct{}, 1)
	go func() {
		<-Wait()
		ch <- struct{}{}
	}()
	start := time.Now()
	startMinute := start.Truncate(1 * time.Minute)
Loop:
	for {
		select {
		case <-ch:
			if startMinute == time.Now().Truncate(1*time.Minute) {
				t.Errorf("stop waiting in the same minute")
			}
			break Loop
		default:
			dur := time.Since(start)
			if dur > 1*time.Minute {
				t.Errorf("waiting more than 1 minute")
			}
		}
	}
}
