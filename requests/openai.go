package requests

import "autflow_back/models"

type CreateAssistantRequest struct {
	Name         string        `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string        `json:"instructions,omitempty" bson:"instructions,omitempty"`
	UserID       string        `json:"user_id" bson:"user_id,omitempty"`
	Collaborator string        `json:"collaborator_id,omitempty" bson:"collaborator_id,omitempty"`
	Type         string        `json:"type" bson:"type"`
	Subs         []models.Subs `json:"subs,omitempty" bson:"subs,omitempty"`
	Active       bool          `json:"active" bson:"active"`
}
