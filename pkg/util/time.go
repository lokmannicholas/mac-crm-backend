package util

import (
	"strconv"
	"strings"
	"time"

	"github.com/lokmannicholas/goplus/timeplus"
	"github.com/lokmannicholas/goplus/timeplus/location"
)

var loc = location.HongKong()

func StrToTimeRange(target string) (time.Time, time.Time, error) {
	dateSplit := strings.Split(target, "-")
	fromStr, err := strconv.ParseInt(dateSplit[0], 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	toStr, err2 := strconv.ParseInt(dateSplit[1], 10, 64)
	if err2 != nil {
		return time.Time{}, time.Time{}, err
	}
	fromTime := time.Unix(0, fromStr*int64(time.Millisecond))
	toTime := time.Unix(0, toStr*int64(time.Millisecond))
	return fromTime, toTime, nil
}

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
