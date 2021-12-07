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
	Version    int    `json:"version"`
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
		res.Items[i] = product{Code: p.Code, PriceCents: p.Price, Version: p.Version}
	}

	c.JSON(http.StatusOK, res)
}

type addProductRequest struct {
	Code       string `json:"code" binding:"required"`
	PriceCents int    `json:"priceCents" binding:"required,min=1"`
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

	c.JSON(http.StatusOK, product{Code: req.Code, Version: 0})
}

type updateProductRequest struct {
	PriceCents int `json:"priceCents" binding:"required,min=1"`
	Version    int `json:"version"`
}

func updateProduct(c *gin.Context) {
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(NewBadRequestError(err))
		return
	}

	db := models.GetDbConnection()
	code := c.Param("code")

	var p models.Product
	if r := db.Where("code = ?", code).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(NewNotFoundError())
		return
	}

	updates := models.Product{
		Price: req.PriceCents,
		BaseModel: models.BaseModel{
			Version: req.Version + 1,
		},
	}

	if r := db.Model(&p).Where("version = ?", req.Version).Updates(updates); r.RowsAffected == 0 {
		c.JSON(NewVersionError())
		return
	}

	c.JSON(http.StatusOK, product{Code: p.Code, PriceCents: p.Price, Version: p.Version})
}

// RegisterProductsRoutes will register all the routes for the products domain
func RegisterProductsRoutes(r *gin.Engine) {
	r.GET("/products", getProducts)
	r.POST("/products", addProduct)
	r.PUT("/products/:code", updateProduct)
}
