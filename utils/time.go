package utils

import "time"

func UnixToStr(unix int64) string {
	u := time.Unix(unix, 0).Format("2006-01-02 15:04:05")
	return u
}