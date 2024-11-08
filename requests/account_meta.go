package requests

import (
	"autflow_back/models"
)

type CreateMetaRequest struct {
	Name          string                `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumberId string                `json:"phone_id" bson:"phone_id"`
	BusinessId    string                `json:"business_id" bson:"business_id"`
	Assistants    []models.AssistantIds `json:"assistants" bson:"assistants"`
	UserID        string                `json:"user_id" bson:"user_id"`
}
