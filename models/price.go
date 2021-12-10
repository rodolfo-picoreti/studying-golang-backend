package models

import (
	"time"

	"github.com/google/uuid"
)

type Price struct {
	Model

	ProductID uuid.UUID
	ChannelID uuid.UUID

	CurrencyCode string `gorm:"not null;index:,sort:desc,type:btree"`
	CentAmount   int    `gorm:"not null"`

	ValidFrom  time.Time `gorm:"not null;index:,sort:desc,type:btree"`
	ValidUntil time.Time `gorm:"index:,sort:desc,type:btree"`
}
