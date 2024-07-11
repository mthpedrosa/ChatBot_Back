package dto

type CreateAssistantDTO struct {
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string `json:"instructions,omitempty" bson:"instructions,omitempty"`
}
