package utils

import "time"

func UnixToStr(unix int64) string {
	u := time.Unix(unix, 0).Format("2006-01-02 15:04:05")
	return u
}

func Today() (start, end time.Time) {
	start, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	end = start.AddDate(0, 0, 1).Add(-1 * time.Second)
	return
}

func ThisMonth() (start, end time.Time) {
	now := time.Now()
	t, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	start = t.AddDate(0, 0, -now.Day()+1)
	end = start.AddDate(0, 1, 0).Add(-1 * time.Second)
	return
}

func ThisWeek() (start, end time.Time) {

	now := time.Now()
	t, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	w := now.Weekday()
	if w == 0 {
		w = 7
	}
	start = t.AddDate(0, 0, -int(w)+1)
	end = start.AddDate(0, 0, 7).Add(-1 * time.Second)
	return
}
