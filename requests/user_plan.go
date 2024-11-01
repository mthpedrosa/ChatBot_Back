package requests

import (
	"autflow_back/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserPlan defines the fields for the user's plan model.
type UserPlanRequest struct {
	UserID       primitive.ObjectID      `bson:"user_id" json:"user_id"`
	PlanType     string                  `bson:"plan_type" json:"plan_type"` // Can be "subscription" or "credit"
	Subscription models.SubscriptionPlan `bson:"subscription,omitempty" json:"subscription,omitempty"`
	Credit       models.CreditPlan       `bson:"credit,omitempty" json:"credit,omitempty"`
}
