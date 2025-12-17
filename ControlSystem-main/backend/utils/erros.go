package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
)

func HandleError(c *gin.Context, err error) {
	if errors.Is(err, ErrUnauthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else if errors.Is(err, ErrForbidden) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	} else if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if errors.Is(err, ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
