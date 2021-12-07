package models

type Product struct {
	BaseModel
	Code  string `gorm:"size:64;unique;not null;index:,sort:desc,type:btree"`
	Price int    `gorm:"not null"`
}
