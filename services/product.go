package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rodolfo-picoreti/studying-golang-backend/models"
	"github.com/rodolfo-picoreti/studying-golang-backend/telemetry"
	"gorm.io/gorm"
)

func getProductCacheKey(sku string) string {
	return fmt.Sprintf("product/%s", sku)
}

func getCache(ctx context.Context, sku string, p *models.Product) error {
	return GetCache(ctx, getProductCacheKey(sku), p)
}

func setCache(ctx context.Context, p *models.Product) error {
	return SetCache(ctx, getProductCacheKey(p.Sku), *p)
}

func expireCache(ctx context.Context, sku string) error {
	return ExpireCache(ctx, getProductCacheKey(sku))
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
		tx = db.Where("sku like ?", fmt.Sprintf("%s%%", codePreffix)).Session(&gorm.Session{})
	}

	products := make([]models.Product, limit)
	tx.Offset(offset).Limit(limit).Find(&products)

	var count int64
	tx.Model(&models.Product{}).Count(&count)

	return &products, count, nil
}

func FindProductBySku(ctx context.Context, sku string) (models.Product, error) {
	_, span := telemetry.Tracer.Start(ctx, "FindProductBySku")
	defer span.End()

	var p models.Product

	if err := getCache(ctx, sku, &p); err != nil {
		db := models.GetDbConnection()

		{
			_, s := telemetry.Tracer.Start(ctx, "SelectDb")
			defer s.End()

			if r := db.Where("sku = ?", sku).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
				return models.Product{}, ProductNotFoundError
			}
		}

		setCache(ctx, &p)
	}

	return p, nil
}

func UpdateProductBySku(ctx context.Context, sku string, updates *models.Product, version int) error {
	_, span := telemetry.Tracer.Start(ctx, "UpdateProductBySku")
	defer span.End()

	db := models.GetDbConnection()

	var p models.Product
	if r := db.Where("sku = ?", sku).First(&p); errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return ProductNotFoundError
	}

	updates.Version = version + 1

	if r := db.Model(&p).Where("version = ?", version).Updates(*updates); r.RowsAffected == 0 {
		return ProductVersionConflictError
	}

	expireCache(ctx, sku)
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
