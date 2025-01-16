// korcen-api/main.go

package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/fluffy-melli/korcen-api/docs"
	"github.com/fluffy-melli/korcen-api/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	fmt.Println("Korcen API Server Start")

	setup := router.SetupRouter()
	setup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := setup.Run(":7777"); err != nil {
		log.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
