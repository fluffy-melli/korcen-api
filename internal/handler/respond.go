// internal/handler/respond.go

package handler

import (
	"github.com/gin-gonic/gin"
)

type ResponseType int

const (
	JSON ResponseType = iota
	XML
)

func respond(c *gin.Context, status int, responseType ResponseType, payload interface{}) {
	switch responseType {
	case XML:
		c.XML(status, payload)
	default:
		c.JSON(status, payload)
	}
}
