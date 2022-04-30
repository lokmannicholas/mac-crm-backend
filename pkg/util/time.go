package util

import (
	"time"

	"github.com/lokmannicholas/goplus/timeplus"
	"github.com/lokmannicholas/goplus/timeplus/location"
)

var loc = location.HongKong()

func TimestampToTime(ts int64) *time.Time {
	if ts == 0 {
		return nil
	}

	tm := time.Unix(ts, 0).In(loc)
	return &tm
}

func StringToTime(s string) *time.Time {
	t, err := timeplus.ParseWithLocation(s, loc)
	if err != nil {
		return nil
	}
	return t
}
func TimeToTimeString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(loc).Format(timeplus.E8601DTwd)
}

func TimeToDateString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(loc).Format(timeplus.E8601DAw)
}
