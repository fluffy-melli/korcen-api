// internal/router/routes.go

package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/fluffy-melli/korcen-api/internal/middleware"
	"github.com/fluffy-melli/korcen-api/pkg/check"
	"github.com/gin-gonic/gin"
)

func respond(c *gin.Context, status int, isXML bool, payload interface{}) {
	if isXML {
		c.XML(status, payload)
	} else {
		c.JSON(status, payload)
	}
}

func SetupRouter(system *actor.ActorSystem, korcenPID *actor.PID, config middleware.MiddlewareConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middleware.TokenBucketMiddleware(config))

	APIGroup := r.Group("/api/v1")
	{
		APIGroup.POST("/korcen", func(c *gin.Context) {
			var header check.Header
			isXML := false

			switch c.ContentType() {
			case "text/xml", "application/xml":
				if err := c.ShouldBindXML(&header); err != nil {
					respond(c, http.StatusBadRequest, true, gin.H{"error": "Invalid XML request"})
					return
				}
				isXML = true
			default:
				if err := c.ShouldBindJSON(&header); err != nil {
					respond(c, http.StatusBadRequest, false, gin.H{"error": "Invalid JSON request"})
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

			future := system.Root.RequestFuture(korcenPID, &check.KorcenRequest{Header: &header}, 5*time.Second)
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
			respond(c, http.StatusOK, isXML, response)
		})
	}

	return r
}
