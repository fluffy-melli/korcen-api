package router

import (
	"net/http"
	"strings"

	_ "github.com/fluffy-melli/korcen-api/docs"
	"github.com/fluffy-melli/korcen-api/internal/middleware"
	"github.com/fluffy-melli/korcen-api/pkg/check"
	"github.com/gin-gonic/gin"
)

// @Summary     Process Korcen Request
// @Description Processes a Korcen request and returns the result
// @Tags        korcen
// @Accept      json,xml
// @Produce     json,xml
// @Param       input  body  check.Header  true  "Korcen Input"
// @Success     200    {object}  check.Respond    "Korcen Result"
// @Failure     400    {object}  map[string]interface{}  "Invalid Request"
// @Router      /api/v1/korcen [post]
func Korcen(c *gin.Context) {
	var header check.Header

	switch c.ContentType() {
	case "text/xml":
		fallthrough
	case "application/xml":
		if err := c.ShouldBindXML(&header); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML request"})
			return
		}
	default:
		if err := c.ShouldBindJSON(&header); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
			return
		}
	}

	if strings.TrimSpace(header.Input) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: empty input"})
		return
	}

	response := check.Korcen(&header)

	if c.GetHeader("Accept") == "application/xml" {
		c.XML(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
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
