package requests

import "autflow_back/models"

type CreateWorkflowRequest struct {
	PhoneMetaId string        `json:"phone_meta_id" bson:"phone_meta_id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Nodes       []models.Node `json:"nodes" bson:"nodes"`
	Active      string        `json:"active" bson:"active"`
	FirstNode   string        `json:"first_node" bson:"first_node"`
	LastNode    string        `json:"last_node" bson:"last_node"`
}

type SendMessageRequest struct {
	PhoneMetaId    string `json:"phone_meta_id" bson:"phone_meta_id"`
	CustomerId     string `json:"customer_id" bson:"customer_id"`
	Message        string `json:"message" bson:"message"`
	ConversationID string `json:"conversation_id" bson:"conversation_id"`
}
