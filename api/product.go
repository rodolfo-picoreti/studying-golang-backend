package api

import (
	"context"
	"errors"
	"example/hello/models"
	"fmt"
	"net/http"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productSearchRequest struct {
	PaginatedRequest
	CodePreffix string `form:"codePreffix"`
}

type product struct {
	Code       string `json:"code"`
	PriceCents int    `json:"priceCents"`
}

func getProducts(c *gin.Context) {
	var req productSearchRequest
	c.BindQuery(&req)

	db := models.GetDbConnection()
	products := make([]models.Product, req.PageLimit)

	tx := db
	if req.CodePreffix != "" {
		tx = db.Where("code like ?", fmt.Sprintf("%s%%", req.CodePreffix))
	}

	var count int64

	if err := crdbgorm.ExecuteTx(context.Background(), tx, nil,
		func(tx *gorm.DB) error {
			tx.Offset(req.Page * req.PageLimit).Limit(req.PageLimit).Find(&products)
			tx.Model(&models.Product{}).Count(&count)
			return nil
		},
	); err != nil {
		fmt.Println(err)
	}

	totalPages := int(count / int64(req.PageLimit))
	if totalPages == 0 {
		totalPages = 1
	}

	res := NewPaginatedResponse(req.Page, totalPages, len(products))
	for i, p := range products {
		res.Items[i] = product{Code: p.Code, PriceCents: p.Price}
	}

	c.JSON(http.StatusOK, res)
}

type addProductRequest struct {
	Code       string `json:"code" binding:"required"`
	PriceCents int    `json:"priceCents" binding:"required,min=1"`
}

type addProductResponse struct {
	Code string `json:"code"`
}

func addProduct(c *gin.Context) {
	var req addProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	db := models.GetDbConnection()

	if r := db.Where("code = ?", req.Code).First(&models.Product{}); !errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(NewAlreadyExistsError())
		return
	}

	db.Create(&models.Product{Code: req.Code, Price: req.PriceCents})

	c.JSON(http.StatusOK, addProductResponse{Code: req.Code})
}

type updateProductRequest struct {
	PriceCents int `json:"priceCents" binding:"required,min=1"`
}

func updateProduct(c *gin.Context) {
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	db := models.GetDbConnection()
	code := c.Param("code")

	var product models.Product
	if r := db.Where("code = ?", code).First(&product); errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(NewNotFoundError())
		return
	}

	product.Price = req.PriceCents
	db.Save(&product)

	c.JSON(http.StatusOK, addProductResponse{Code: code})
}

// RegisterProductsRoutes will register all the routes for the products domain
func RegisterProductsRoutes(r *gin.Engine) {
	r.GET("/products", getProducts)
	r.POST("/products", addProduct)
	r.PUT("/products/:code", updateProduct)
}
