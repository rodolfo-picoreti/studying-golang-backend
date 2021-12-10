package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Product struct {
	Model

	Sku         string `gorm:"size:64;unique;not null;index:,sort:desc,type:btree"`
	Name        string `gorm:"not null;index:,sort:desc,type:btree"`
	Description string
	Status      string         `gorm:"not null;index:,sort:desc,type:btree"`
	Tags        pq.StringArray `gorm:"type:text[]"`

	ProductTypeID uuid.UUID
	ProductType   ProductType

	Attributes  []Attribute
	Inventories []Inventory
	Prices      []Price
}
