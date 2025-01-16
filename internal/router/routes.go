package router

import (
	"strings"

	_ "github.com/fluffy-melli/korcen-api/docs"
	"github.com/fluffy-melli/korcen-api/internal/middleware"
	"github.com/fluffy-melli/korcen-api/pkg/check"
	"github.com/gin-gonic/gin"
)

// @Summary		Process Korcen Request
// @Description	Processes a Korcen request and returns the result
// @Tags			korcen
// @Accept			json
// @Produce		json
// @Param			input	body		check.Header	true	"Korcen Input"
// @Success		200		{object}	check.Respond	"Korcen Result"
// @Failure		400		{object}	map[string]interface{}	"Invalid Request"
// @Router			/api/v1/korcen [post]
func Korcen(c *gin.Context) {
	var header check.Header

	if err := c.ShouldBindJSON(&header); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if strings.TrimSpace(header.Input) == "" {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	response := check.Korcen(&header)

	c.JSON(200, response)
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.TokenBucketMiddleware())
	APIGroup := r.Group("/api/v1")
	{
		APIGroup.POST("/korcen", Korcen)
	}

	return r
}
