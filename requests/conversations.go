package requests

import "autflow_back/models"

type ConversationRequest struct {
	CustomerId  string           `json:"customer_id" bson:"customer_id"`
	Messages    []models.Message `json:"mensagens" bson:"mensagens,omitempty"`
	WorkflowID  string           `json:"workflow_id" bson:"workflow_id"`
	OtherFields []models.Fields  `json:"other_fields" bson:"other_fields"`
}
