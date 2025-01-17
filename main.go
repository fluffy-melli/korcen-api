// main.go

package main

import (
	"fmt"
	"log"
	"os"

	// Proto.Actor
	"github.com/asynkron/protoactor-go/actor"

	_ "github.com/fluffy-melli/korcen-api/docs"
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

	r := router.SetupRouter(system, korcenPID)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":7777"); err != nil {
		log.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
