package services

import (
	"context"
	"errors"
	"example/hello/models"
	"example/hello/telemetry"
	"fmt"

	"gorm.io/gorm"
)

func getProductCacheKey(code string) string {
	return fmt.Sprintf("product/%s", code)
}

func getCache(ctx context.Context, code string, p *models.Product) error {
	return GetCache(ctx, getProductCacheKey(code), p)
}

func setCache(ctx context.Context, p *models.Product) error {
	return SetCache(ctx, getProductCacheKey(p.Code), *p)
}

func expireCache(ctx context.Context, code string) error {
	return ExpireCache(ctx, getProductCacheKey(code))
}

var (
	ProductNotFoundError        = errors.New("Product not found")
	ProductVersionConflictError = errors.New("Product version not found")
	ProductCreateError          = errors.New("Product creation failed")
)

func FindProducts(ctx context.Context, offset int, limit int, codePreffix string) (*[]models.Product, int64, error) {
	_, span := telemetry.Tracer.Start(ctx, "FindProducts")
	defer span.End()

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

func FindProductByCode(ctx context.Context, code string) (models.Product, error) {
	_, span := telemetry.Tracer.Start(ctx, "FindProductByCode")
	defer span.End()

	var p models.Product

	if err := getCache(ctx, code, &p); err != nil {
		db := models.GetDbConnection()

		{
			_, s := telemetry.Tracer.Start(ctx, "SelectDb")
			defer s.End()

			if r := db.Where("code = ?", code).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
				return models.Product{}, ProductNotFoundError
			}
		}

		setCache(ctx, &p)
	}

	return p, nil
}

func UpdateProductByCode(ctx context.Context, code string, updates *models.Product, version int) error {
	_, span := telemetry.Tracer.Start(ctx, "UpdateProductByCode")
	defer span.End()

	db := models.GetDbConnection()

	var p models.Product
	if r := db.Where("code = ?", code).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return ProductNotFoundError
	}

	updates.Version = version + 1

	if r := db.Model(&p).Where("version = ?", version).Updates(*updates); r.RowsAffected == 0 {
		return ProductVersionConflictError
	}

	expireCache(ctx, code)
	return nil
}

func CreateProduct(ctx context.Context, p *models.Product) error {
	_, span := telemetry.Tracer.Start(ctx, "CreateProduct")
	defer span.End()

	db := models.GetDbConnection()

	if r := db.Create(p); r.Error != nil {
		return ProductCreateError
	}

	return nil
}
