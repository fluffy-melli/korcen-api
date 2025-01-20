package handler

import "github.com/gin-gonic/gin"

func respond(c *gin.Context, status int, isXML bool, payload interface{}) {
	if isXML {
		c.XML(status, payload)
	} else {
		c.JSON(status, payload)
	}
}
