package requests

type CreateAssistantRequest struct {
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string `json:"instructions,omitempty" bson:"instructions,omitempty"`
	UserID       string `json:"user_id,omitempty" bson:"instructions,omitempty"`
}
