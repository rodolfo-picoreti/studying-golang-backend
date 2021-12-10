package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodolfo-picoreti/studying-golang-backend/models"
	"github.com/rodolfo-picoreti/studying-golang-backend/services"
)

type createProductTypeRequest struct {
}

func createProductType(c *gin.Context) {
	ctx := c.Request.Context()
	var req createProductTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	_, err := services.FindProductBySku(ctx, req.Sku)
	if err == nil {
		c.JSON(NewAlreadyExistsError())
		return
	}

	// todo: build model
	services.CreateProduct(ctx, &models.Product{})

	c.Status(http.StatusOK)
}

// RegisterProductsRoutes will register all the routes for the products domain
func RegisterProductTypesRoutes(r *gin.Engine) {
	r.POST("/product-types", createProductType)
}
