package models

import (
	"time"
)

// Config represents a configuration entity
type Config struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Value       string    `bson:"value" json:"value"`
	Type        string    `bson:"type" json:"type"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
