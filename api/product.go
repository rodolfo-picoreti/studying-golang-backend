package api

import (
	"example/hello/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productSearchRequest struct {
	Page      int `form:"page"`
	PageLimit int `form:"pageLimit"`

	CodePreffix string `form:"codePreffix"`
}

type product struct {
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

func getProducts(c *gin.Context) {
	var req productSearchRequest
	c.BindQuery(&req)

	db := models.GetDbConnection()
	products := make([]models.Product, req.PageLimit)

	db.Offset(req.Page * req.PageLimit).Limit(req.PageLimit).Find(&products)

	res := NewPaginatedItems(0, 0, len(products))
	for i, p := range products {
		res.Items[i] = product{Code: p.Code, Price: p.Price}
	}

	c.JSON(http.StatusOK, res)
}

type addProductRequest struct {
	Code string `json:"code" binding:"required"`
}

type addProductResponse struct {
	Code string `json:"code"`
}

func addProduct(c *gin.Context) {
	var req addProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := models.GetDbConnection()

	product := models.Product{Code: req.Code}
	db.Create(&product)

	c.JSON(http.StatusOK, addProductResponse{Code: product.Code})
}

// RegisterProductsRoutes will register all the routes for the products domain
func RegisterProductsRoutes(r *gin.Engine) {
	r.GET("/products", getProducts)
	r.POST("/products", addProduct)
}
