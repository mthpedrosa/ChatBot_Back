package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserPlan defines the fields for the user's plan model.
type UserPlan struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	PlanType     string             `bson:"plan_type" json:"plan_type"` // Can be "subscription" or "credit"
	Subscription *SubscriptionPlan  `bson:"subscription,omitempty" json:"subscription,omitempty"`
	Credit       *CreditPlan        `bson:"credit,omitempty" json:"credit,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// SubscriptionPlan defines the fields for subscription plans with limited messages.
type SubscriptionPlan struct {
	MessagesRemaining int       `bson:"messages_remaining" json:"messages_remaining"`
	TotalMessages     int       `bson:"total_messages" json:"total_messages"`
	ExpirationDate    time.Time `bson:"expiration_date" json:"expiration_date"`
}

// CreditPlan defines the fields for credit-based plans, where each message deducts a value from the balance.
type CreditPlan struct {
	Balance        float64 `bson:"balance" json:"balance"`                   // Balance in credits
	CostPerMessage float64 `bson:"cost_per_message" json:"cost_per_message"` // Cost of each message in credits
}
