package dto

import (
	"autflow_back/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationCreateDTO struct {
	CustomerId  string           `json:"customer_id" bson:"customer_id"`
	Messages    []models.Message `json:"mensagens" bson:"mensagens,omitempty"`
	WorkflowID  string           `json:"workflow_id" bson:"workflow_id"`
	OtherFields []models.Fields  `json:"other_fields" bson:"other_fields"`
}

type ConversationListDTO struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerId string             `json:"customer_id" bson:"customer_id"`
	WorkflowID string             `json:"workflow_id" bson:"workflow_id"`
}

type ConversationDetailDTO struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerId  string             `json:"customer_id" bson:"customer_id"`
	Messages    []models.Message   `json:"mensagens" bson:"mensagens,omitempty"`
	WorkflowID  string             `json:"workflow_id" bson:"workflow_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	OtherFields []models.Fields    `json:"other_fields" bson:"other_fields"`
}

func (dto *ConversationCreateDTO) ToConversation() models.Conversation {
	return models.Conversation{
		CustomerId:  dto.CustomerId,
		Messages:    dto.Messages,
		WorkflowID:  dto.WorkflowID,
		OtherFields: dto.OtherFields,
	}
}
