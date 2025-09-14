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

func ParseYYYYMMDD(s string) (y int, m time.Month, d int, err error) {
	t, e := time.Parse("2006-01-02", s)
	if e != nil {
		return 0, 0, 0, e
	}
	return t.Year(), t.Month(), t.Day(), nil
}

func ParseCutoffHHMMSS(s string) (h, m, sec int, err error) {
	if tm, e := time.Parse("15:04:05", s); e == nil {
		return tm.Hour(), tm.Minute(), tm.Second(), nil
	}
	if tm, e := time.Parse("15:04", s); e == nil {
		return tm.Hour(), tm.Minute(), 0, nil
	}
	return 0, 0, 0, fmt.Errorf("invalid cutoff time %q", s)
}

func LoadLocationOrUTC(tz string) *time.Location {
	if tz == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}

func DayBoundsLocalToUTC(loc *time.Location, y int, m time.Month, d int) (fromUTC, toUTC time.Time) {
	startLocal := time.Date(y, m, d, 0, 0, 0, 0, loc)
	endLocal := time.Date(y, m, d, 23, 59, 59, int(time.Nanosecond*999999999), loc)
	return startLocal.UTC(), endLocal.UTC()
}

type Pagination struct {
    TotalData   int64 `json:"totalData"`
    CurrentPage int   `json:"currentPage"`
    TotalPages  int   `json:"totalPages"`
    HasNextPage bool  `json:"hasNextPage"`
    HasPrevPage bool  `json:"hasPrevPage"`
}

func BuildPagination(total int64, page, limit int) Pagination {
    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }

    totalPages := int((total + int64(limit) - 1) / int64(limit))
    if totalPages == 0 {
        totalPages = 1
    }
    hasNext := page < totalPages
    hasPrev := page > 1

    return Pagination{
        TotalData:   total,
        CurrentPage: page,
        TotalPages:  totalPages,
        HasNextPage: hasNext,
        HasPrevPage: hasPrev,
    }
}