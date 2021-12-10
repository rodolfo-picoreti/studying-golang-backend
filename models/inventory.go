package models

import (
	"github.com/google/uuid"
)

type Inventory struct {
	Model

	ProductID uuid.UUID
	ChannelID uuid.UUID

	AvailableQuantity int
}
