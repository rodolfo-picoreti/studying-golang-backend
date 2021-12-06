package models

type Product struct {
	BaseModel
	Code  string `gorm:"index:idx_code,unique"`
	Price uint   `gorm:"not_null,check:price>0"`
}
