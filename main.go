// main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asynkron/protoactor-go/actor"

	_ "github.com/fluffy-melli/korcen-api/docs"
	"github.com/fluffy-melli/korcen-api/internal/middleware"
	"github.com/fluffy-melli/korcen-api/internal/router"
	"github.com/fluffy-melli/korcen-api/pkg/check"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	fmt.Println("Korcen API Server Start")

	system := actor.NewActorSystem()

	props := actor.PropsFromProducer(func() actor.Actor {
		return &check.KorcenActor{}
	})

	korcenPID := system.Root.Spawn(props)

	config := middleware.MiddlewareConfig{
		Capacity:   100,          // 최대 토큰 수
		RefillRate: 100.0 / 60.0, // 초당 리필 속도 (1분당 100개)
	}

	r := router.SetupRouter(system, korcenPID, config)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":7777"); err != nil {
		log.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
