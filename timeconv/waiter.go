package timeconv

import "time"

func Wait() <-chan time.Time {
	t := time.Now()
	t = t.Truncate(1 * time.Minute)
	t = t.Add(1 * time.Minute)
	return time.After(time.Until(t))
}
