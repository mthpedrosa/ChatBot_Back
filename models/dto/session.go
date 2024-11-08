package dto

import (
	"autflow_back/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionCreateDTO struct {
	CustomerID     string           `json:"customer_id,omitempty" bson:"customer_id"`
	AssistantId    string           `json:"assistant_id,omitempty" bson:"assistant_id"`
	ConversationId string           `json:"conversation_id" bson:"conversation_id"`
	Status         string           `json:"status" bson:"status"`
	Tags           []string         `json:"tags,omitempty" bson:"tags,omitempty"`
	OtherFields    []models.Fields  `json:"other_fields,omitempty" bson:"other_fields"`
	Messages       []models.Message `json:"messages,omitempty" bson:"messages,omitempty"`
	LastNode       string           `json:"last_node,omitempty" bson:"last_node"`
}

type SessionListDTO struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID     string             `json:"customer_id,omitempty" bson:"customer_id"`
	AssistantId    string             `json:"assistant_id,omitempty" bson:"assistant_id"`
	ConversationId string             `json:"conversation_id" bson:"conversation_id"`
	Status         string             `json:"status" bson:"status"`
}

type SessionDetailDTO struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID     string             `json:"customer_id,omitempty" bson:"customer_id"`
	AssistantId    string             `json:"assistant_id,omitempty" bson:"assistant_id"`
	ConversationId string             `json:"conversation_id" bson:"conversation_id"`
	Status         string             `json:"status" bson:"status"`
	Duration       string             `json:"duration" bson:"duration"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt       time.Time          `json:"update_at" bson:"update_at"`
	FinishedAt     time.Time          `json:"finished_at" bson:"finished_at"`
	Tags           []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	OtherFields    []models.Fields    `json:"other_fields,omitempty" bson:"other_fields"`
	Messages       []models.Message   `json:"messages,omitempty" bson:"messages,omitempty"`
	LastNode       string             `json:"last_node,omitempty" bson:"last_node"`
}

func (dto *SessionCreateDTO) ToSession() models.Session {
	return models.Session{
		CustomerID:     dto.CustomerID,
		AssistantId:    dto.AssistantId,
		ConversationId: dto.ConversationId,
		Status:         dto.Status,
		Tags:           dto.Tags,
		OtherFields:    dto.OtherFields,
		Messages:       dto.Messages,
		LastNode:       dto.LastNode,
	}
}

func (dto *SessionCreateDTO) Validate() error {
	if dto.CustomerID == "" {
		return errors.New("O CustomerID é obrigatório e não pode estar em branco.")
	}
	if dto.ConversationId == "" {
		return errors.New("O ConversationId é obrigatório e não pode estar em branco.")
	}

	if dto.AssistantId == "" {
		return errors.New("O AssistantId é obrigatório e não pode estar em branco.")
	}

	return nil
}
