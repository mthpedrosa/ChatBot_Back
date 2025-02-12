package dto

import (
	"autflow_back/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssistantCreateDTO struct {
	Name         string        `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string        `json:"instructions,omitempty" bson:"instructions,omitempty"`
	UserID       string        `json:"user_id" bson:"user_id,omitempty"`
	Collaborator string        `json:"collaborator_id,omitempty" bson:"collaborator_id,omitempty"`
	Type         string        `json:"type" bson:"type"`
	Subs         []models.Subs `json:"subs,omitempty" bson:"subs,omitempty"`
	Active       bool          `json:"active" bson:"active"`
	Info         string        `json:"info,omitempty" bson:"info,omitempty"`
}

type AssitantListDTO struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Active bool               `json:"active" bson:"active"`
	Type   string             `json:"type" bson:"type"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
}

type AssitantDetailDTO struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string             `json:"instructions,omitempty" bson:"instructions,omitempty"`
	UserID       string             `json:"user_id" bson:"user_id,omitempty"`
	Collaborator string             `json:"collaborator_id,omitempty" bson:"collaborator_id,omitempty"`
	Type         string             `json:"type" bson:"type"`
	Subs         []models.Subs      `json:"subs" bson:"subs"`
	Active       bool               `json:"active" bson:"active"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt     time.Time          `json:"update_at" bson:"update_at"`
	IdAssistant  string             `json:"assistant_id,omitempty" bson:"assistant_id,omitempty"`
}

func (dto *AssistantCreateDTO) ToAssistant() models.CreateAssistant {
	return models.CreateAssistant{
		Name:         dto.Name,
		Instructions: dto.Instructions,
		UserID:       dto.UserID,
		Collaborator: dto.Collaborator,
		Subs:         dto.Subs,
		Active:       dto.Active,
		Type:         dto.Type,
		Info:         dto.Info,
	}
}
