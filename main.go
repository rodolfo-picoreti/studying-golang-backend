package main

import (
	"example/hello/api"
	"example/hello/models"
	"example/hello/telemetry"

	"github.com/gin-gonic/gin"
)

func main() {
	telemetry.InitLogger()

	shutdown := telemetry.InitTraceProvider()
	defer shutdown()

	models.InitDB()
	models.AutoMigrate()

	r := gin.New()
	r.Use(telemetry.TraceMiddleware())
	r.Use(telemetry.LoggerMiddleware())

	api.RegisterProductsRoutes(r)
	r.Run(":8080")
}
