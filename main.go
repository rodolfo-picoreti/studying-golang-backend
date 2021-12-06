package main

import (
	"example/hello/api"

	"github.com/gin-gonic/gin"
)

func main() {
	// db := models.GetDbConnection()
	// db.AutoMigrate(&models.Product{})

	r := gin.Default()

	api.RegisterProductsRoutes(r)
	r.Run()
}
