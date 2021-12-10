package models

import "github.com/google/uuid"

type AttributeDefinition struct {
	Model

	ProductTypeID uuid.UUID

	Type         string `gorm:"not null;index:,sort:desc,type:btree"`
	Name         string `gorm:"not null;index:,sort:desc,type:btree"`
	Description  string
	IsSearchable bool
	Validation   string
}
