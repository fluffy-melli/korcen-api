// hi/korcen-api/main.go

package main

import (
	_ "github.com/fluffy-melli/korcen-api/docs"
	"github.com/fluffy-melli/korcen-api/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	setup := router.SetupRouter()
	setup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	setup.Run(":7777")
}
