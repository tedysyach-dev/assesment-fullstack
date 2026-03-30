package utils

import (
	"time"
)

func FormatTime(t *time.Time, timeType string) string {
	if t == nil {
		return ""
	}
	if t.IsZero() {
		return ""
	}

	if timeType == "date" {
		return t.Format("2006-01-02")
	} else {
		return t.Format("2006-01-02 15:04:05")
	}
}
