package helper

import (
	"errors"
	"fmt"
	"time"
)

func NormalizeStringField(name string) string {
	return string([]rune(name))
}

func ParseTimeOfDay(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("empty time")
	}
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("time must be HH:MM:SS, got %q", s)
	}

	return time.Date(1, time.January, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC), nil
}
