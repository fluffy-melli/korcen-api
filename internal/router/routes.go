// internal/router/routes.go
package router

import (
	"net/http"
	"strings"
	"time"

	_ "github.com/fluffy-melli/korcen-api/docs"
	"github.com/fluffy-melli/korcen-api/internal/middleware"
	"github.com/fluffy-melli/korcen-api/pkg/check"
	"github.com/gin-gonic/gin"

	"github.com/asynkron/protoactor-go/actor"
)

func SetupRouter(system *actor.ActorSystem, korcenPID *actor.PID) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(middleware.TokenBucketMiddleware())

	APIGroup := r.Group("/api/v1")
	{
		APIGroup.POST("/korcen", func(c *gin.Context) {
			var header check.Header
			isXML := false

			switch c.ContentType() {
			case "text/xml", "application/xml":
				if err := c.ShouldBindXML(&header); err != nil {
					c.XML(http.StatusBadRequest, gin.H{"error": "Invalid XML request"})
					return
				}
				isXML = true
			default:
				if err := c.ShouldBindJSON(&header); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
					return
				}
			}

			if strings.TrimSpace(header.Input) == "" {
				if isXML {
					c.XML(http.StatusBadRequest, gin.H{"error": "Invalid request: empty input"})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: empty input"})
				}
				return
			}

			msg := &check.KorcenRequest{Header: &header}
			future := system.Root.RequestFuture(korcenPID, msg, 5*time.Second)

			result, err := future.Result()
			if err != nil {
				if isXML {
					c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return
			}

			korcenResp, ok := result.(*check.KorcenResponse)
			if !ok {
				if isXML {
					c.XML(http.StatusInternalServerError, gin.H{"error": "Invalid actor response"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid actor response"})
				}
				return
			}

			if korcenResp.Err != nil {
				if isXML {
					c.XML(http.StatusBadRequest, gin.H{"error": korcenResp.Err.Error()})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": korcenResp.Err.Error()})
				}
				return
			}

			response := korcenResp.Respond

			if isXML {
				c.XML(http.StatusOK, response)
			} else {
				c.JSON(http.StatusOK, response)
			}
		})
	}

	return r
}
