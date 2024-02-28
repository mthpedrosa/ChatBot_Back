package requests

type CreateUserRequest struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Profile  string `json:"profile,omitempty" bson:"profile,omitempty"`
}
