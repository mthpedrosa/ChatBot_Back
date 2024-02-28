package requests

import "autflow_back/models"

type CreateMetaRequest struct {
	Name           string              `json:"name,omitempty" bson:"name,omitempty"`
	Token          string              `json:"token" bson:"token"`
	MetaID         string              `json:"meta_id" bson:"meta_id"`
	PhonesMeta     []models.PhonesMeta `json:"phones_meta" bson:"phones_meta"`
	EditPermission []string            `json:"edit_permission" bson:"edit_permission"`
}
