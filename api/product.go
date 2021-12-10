package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodolfo-picoreti/studying-golang-backend/models"
	"github.com/rodolfo-picoreti/studying-golang-backend/services"
)

type product struct {
	Sku         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     int    `json:"version"`
}

type productSearchRequest struct {
	PaginatedRequest
	CodePreffix string `form:"codePreffix"`
}

type updateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     int    `json:"version"`
}

type addProductRequest struct {
	Sku         string `json:"sku" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func buildProductDto(p *models.Product) product {
	return product{
		Sku:     p.Sku,
		Version: p.Version,
	}
}

func getProducts(c *gin.Context) {
	ctx := c.Request.Context()
	var req productSearchRequest
	c.BindQuery(&req)

	offset := req.Page * req.PageLimit
	limit := req.PageLimit
	codePreffix := req.CodePreffix

	products, count, _ := services.FindProducts(ctx, offset, limit, codePreffix)

	totalPages := int(count / int64(req.PageLimit))
	if totalPages == 0 {
		totalPages = 1
	}

	res := NewPaginatedResponse(req.Page, totalPages, len(*products))
	for i, p := range *products {
		res.Items[i] = buildProductDto(&p)
	}

	c.JSON(http.StatusOK, res)
}

func getProduct(c *gin.Context) {
	ctx := c.Request.Context()
	sku := c.Param("sku")

	p, err := services.FindProductBySku(ctx, sku)
	if err != nil {
		c.JSON(NewNotFoundError())
		return
	}

	c.JSON(http.StatusOK, buildProductDto(&p))
}

func createProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req addProductRequest
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

func updateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	sku := c.Param("sku")
	// todo: build model
	updates := models.Product{}

	if err := services.UpdateProductBySku(ctx, sku, &updates, req.Version); err != nil {
		switch err {
		case services.ProductNotFoundError:
			c.JSON(NewNotFoundError())
		case services.ProductVersionConflictError:
			c.JSON(NewVersionError())
		}
	}

	c.Status(http.StatusOK)
}

// RegisterProductsRoutes will register all the routes for the products domain
func RegisterProductsRoutes(r *gin.Engine) {
	r.GET("/products", getProducts)
	r.GET("/products/:sku", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:sku", updateProduct)
}
