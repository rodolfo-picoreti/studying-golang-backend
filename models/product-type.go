package models

type ProductType struct {
	Model

	Name        string `gorm:"not null;index:,sort:desc,type:btree"`
	Description string

	AttributesDefinitions []AttributeDefinition
}
