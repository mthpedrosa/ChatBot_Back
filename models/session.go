package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID     string             `json:"customer_id,omitempty" bson:"customer_id"`
	WorkflowId     string             `json:"workflow_id,omitempty" bson:"workflow_id"`
	ConversationId string             `json:"conversation_id" bson:"conversation_id"`
	Status         string             `json:"status" bson:"status"`
	Duration       string             `json:"duration" bson:"duration"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt       time.Time          `json:"update_at" bson:"update_at"`
	FinishedAt     time.Time          `json:"finished_at" bson:"finished_at"`
	Tags           []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	OtherFields    []Fields           `json:"other_fields,omitempty" bson:"other_fields"`
	Messages       []Message          `json:"messages,omitempty" bson:"messages,omitempty"`
	LastNode       string             `json:"last_node,omitempty" bson:"last_node"`
}
