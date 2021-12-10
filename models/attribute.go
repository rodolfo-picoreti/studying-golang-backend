package models

import "github.com/google/uuid"

type Attribute struct {
	Model
	ProductID uuid.UUID
	Key       string
	Value     string
}
