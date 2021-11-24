package timeutil

import (
	"time"
)

type Time interface {
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Now() time.Time
	Since(t time.Time) time.Duration
}

func NewTime() Time {
	return &realTime{}
}

type realTime struct{}

func (rt *realTime) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (rt *realTime) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (rt *realTime) Now() time.Time {
	return time.Now()
}

func (rt *realTime) Since(t time.Time) time.Duration {
	return rt.Now().Sub(t)
}

func ConvertToTimeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func ConvertAndValidateTimeToMillis(time int64) int64 {
	if time < 9999999999 {
		return time * 1000
	}
	return time
}

func ConvertMinuteToMillis(minute int64) int64 {
	return (minute * int64(time.Minute)) / int64(time.Millisecond)
}

func ConvertTimeToMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
