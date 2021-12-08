package api

import (
	"example/hello/models"
	"example/hello/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type product struct {
	Code       string `json:"code"`
	PriceCents int    `json:"priceCents"`
	Version    int    `json:"version"`
}

type productSearchRequest struct {
	PaginatedRequest
	CodePreffix string `form:"codePreffix"`
}

type updateProductRequest struct {
	PriceCents int `json:"priceCents" binding:"required,min=1"`
	Version    int `json:"version"`
}

type addProductRequest struct {
	Code       string `json:"code" binding:"required"`
	PriceCents int    `json:"priceCents" binding:"required,min=1"`
}

func getProducts(c *gin.Context) {
	var req productSearchRequest
	c.BindQuery(&req)

	offset := req.Page * req.PageLimit
	limit := req.PageLimit
	codePreffix := req.CodePreffix

	products, count, _ := services.FindProducts(offset, limit, codePreffix)

	totalPages := int(count / int64(req.PageLimit))
	if totalPages == 0 {
		totalPages = 1
	}

	res := NewPaginatedResponse(req.Page, totalPages, len(*products))
	for i, p := range *products {
		res.Items[i] = product{Code: p.Code, PriceCents: p.Price, Version: p.Version}
	}

	c.JSON(http.StatusOK, res)
}

func getProduct(c *gin.Context) {
	code := c.Param("code")

	p, err := services.FindProductByCode(code)
	if err != nil {
		c.JSON(NewNotFoundError())
		return
	}

	c.JSON(http.StatusOK, product{Code: p.Code, PriceCents: p.Price})
}

func createProduct(c *gin.Context) {
	var req addProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	_, err := services.FindProductByCode(req.Code)
	if err != nil {
		c.JSON(NewNotFoundError())
		return
	}

	services.CreateProduct(&models.Product{Code: req.Code, Price: req.PriceCents})

	c.Status(http.StatusOK)
}

func updateProduct(c *gin.Context) {
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	code := c.Param("code")
	updates := models.Product{
		Price: req.PriceCents,
	}

	if err := services.UpdateProductByCode(code, &updates, req.Version); err != nil {
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
	r.GET("/products/:code", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:code", updateProduct)
}
