package services

import (
	"errors"
	"example/hello/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func getProductCacheKey(code string) string {
	return fmt.Sprintf("product/%s", code)
}

var (
	ProductNotFoundError        = errors.New("Product not found")
	ProductVersionConflictError = errors.New("Product version not found")
	ProductCreateError          = errors.New("Product creation failed")
)

func FindProducts(offset int, limit int, codePreffix string) (*[]models.Product, int64, error) {
	db := models.GetDbConnection()

	tx := db
	if codePreffix != "" {
		tx = db.Where("code like ?", fmt.Sprintf("%s%%", codePreffix)).Session(&gorm.Session{})
	}

	products := make([]models.Product, limit)
	tx.Offset(offset).Limit(limit).Find(&products)

	var count int64
	tx.Model(&models.Product{}).Count(&count)

	return &products, count, nil
}

func FindProductByCode(code string) (models.Product, error) {
	var p models.Product

	if err := GetCacheStore().Get(getProductCacheKey(code), &p); err != nil {
		db := models.GetDbConnection()

		if r := db.Where("code = ?", code).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return models.Product{}, ProductNotFoundError
		}
	}

	GetCacheStore().Set(getProductCacheKey(code), p, time.Minute)
	return p, nil
}

func UpdateProductByCode(code string, updates *models.Product, version int) error {
	db := models.GetDbConnection()

	var p models.Product
	if r := db.Where("code = ?", code).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return ProductNotFoundError
	}

	updates.Version = version + 1

	if r := db.Model(&p).Where("version = ?", version).Updates(*updates); r.RowsAffected == 0 {
		return ProductVersionConflictError
	}

	GetCacheStore().Delete(getProductCacheKey(code))
	return nil
}

func CreateProduct(product *models.Product) error {
	db := models.GetDbConnection()

	if r := db.Create(product); r.Error != nil {
		return ProductCreateError
	}

	GetCacheStore().Set(getProductCacheKey(product.Code), *product, time.Minute)

	return nil
}
