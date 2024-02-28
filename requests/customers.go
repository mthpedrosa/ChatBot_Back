package requests

import "autflow_back/models"

type CustomerRequest struct {
	Name        string          `json:"name,omitempty" bson:"name,omitempty"`
	Email       string          `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string          `json:"phone,omitempty" bson:"phone,omitempty"`
	WhatsAppID  string          `json:"whatsapp_id,omitempty" bson:"whatsapp_id,omitempty"`
	OtherFields []models.Fields `json:"other_fields" bson:"other_fields"`
}
