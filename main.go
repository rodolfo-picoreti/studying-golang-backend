package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rodolfo-picoreti/studying-golang-backend/api"
	"github.com/rodolfo-picoreti/studying-golang-backend/models"
	"github.com/rodolfo-picoreti/studying-golang-backend/telemetry"
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
