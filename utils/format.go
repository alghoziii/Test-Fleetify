package utils

import (
	"time"
)

const TimeLayout = "2006-01-02 15:04"

func FormatTS(t time.Time) string {
	return t.Format(TimeLayout)
}

func FormatTSPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(TimeLayout)
	return &s
}
