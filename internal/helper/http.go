package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsaFan/fleetify-be/internal/appErr"
)

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

func NotFound(c *gin.Context, msg string) {
	RespondErr(c, http.StatusNotFound, "not_found", msg)
}

func WriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, appErr.ErrAlreadyExists):
		Conflict(c, err.Error())
	case errors.Is(err, appErr.ErrRequiredField),
		errors.Is(err, appErr.ErrInvalidInput),
		errors.Is(err, appErr.ErrInvalidRange),
		errors.Is(err, appErr.ErrInvalidTimeRange):
		BadRequest(c, err.Error())
	case errors.Is(err, appErr.ErrNotFound):
		NotFound(c, err.Error())
	default:
		Internal(c, err.Error())
	}
}
