// internal/router/routes.go

package router

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/fluffy-melli/korcen-api/internal/handler"
	"github.com/fluffy-melli/korcen-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(system *actor.ActorSystem, korcenPID *actor.PID, config middleware.MiddlewareConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middleware.TokenBucketMiddleware(config))

	APIGroup := r.Group("/api/v1")
	{
		APIGroup.POST("/korcen", func(c *gin.Context) {
			handler.KorcenV1(c, system, korcenPID)
		})
	}

	return r
}
