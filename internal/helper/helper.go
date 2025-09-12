package helper

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func RespondErr(c *gin.Context, code int, err string, msg string) {
	c.JSON(code, gin.H{
		"error":   err,
		"message": msg,
	})
	c.Abort()
}

func BadRequest(c *gin.Context, msg string) {
	RespondErr(c, http.StatusBadRequest, "bad_request", msg)
}

func Conflict(c *gin.Context, msg string) {
	RespondErr(c, http.StatusConflict, "conflict", msg)
}

func Internal(c *gin.Context, msg string) {
	RespondErr(c, http.StatusInternalServerError, "internal_error", msg)
}
