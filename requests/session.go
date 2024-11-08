package requests

import "autflow_back/models"

type SessionRequest struct {
	CustomerID     string           `json:"customer_id,omitempty" bson:"customer_id"`
	AssistantId    string           `json:"assistant_id,omitempty" bson:"assistant_id"`
	ConversationId string           `json:"conversation_id" bson:"conversation_id"`
	Status         string           `json:"status" bson:"status"`
	Tags           []string         `json:"tags,omitempty" bson:"tags,omitempty"`
	OtherFields    []models.Fields  `json:"other_fields,omitempty" bson:"other_fields"`
	Messages       []models.Message `json:"messages,omitempty" bson:"messages,omitempty"`
	LastNode       string           `json:"last_node,omitempty" bson:"last_node"`
}

type SessionOtherRequest struct {
	OtherFields models.Fields `json:"other_fields,omitempty" bson:"other_fields"`
}
