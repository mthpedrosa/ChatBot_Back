package requests

type CreateAssistantRequest struct {
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string `json:"instructions,omitempty" bson:"instructions,omitempty"`
	IdCustomer   string `json:"customerid,omitempty" bson:"instructions,omitempty"`
}
