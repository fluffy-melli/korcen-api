// main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		RefillRate: 100.0 / 60.0, // 초당 리필 속도 (1분당 100회)
	}

	r := router.SetupRouter(system, korcenPID, config)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":7777",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()
	log.Println("Server is running on port 7777.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Shutting down worker pool...")
	check.ShutdownWorkerPool()
	log.Println("Worker pool has been shut down successfully.")

	log.Println("Server exiting.")
}
