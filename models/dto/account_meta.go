package dto

import (
	"autflow_back/models"
	"autflow_back/requests"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateMetaDTO struct {
	Name          string                `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumberId string                `json:"phone_id" bson:"phone_id"`
	BusinessId    string                `json:"business_id" bson:"business_id"`
	Assistants    []models.AssistantIds `json:"assistants" bson:"assistants"`
	UserID        string                `json:"user_id" bson:"user_id"`
}
type MetaListDTO struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumberId string             `json:"phone_id" bson:"phone_id"`
	BusinessId    string             `json:"business_id" bson:"business_id"`
	UserID        string             `json:"user_id" bson:"user_id"`
}

type MetaDetailDTO struct {
	ID            primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string                `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumberId string                `json:"phone_id" bson:"phone_id"`
	BusinessId    string                `json:"business_id" bson:"business_id"`
	CreatedAt     time.Time             `json:"created_at" bson:"created_at"`
	UpdateAt      time.Time             `json:"update_at" bson:"update_at"`
	Assistants    []models.AssistantIds `json:"assistants" bson:"assistants"`
	UserID        string                `json:"user_id" bson:"user_id"`
}

func NewCreateMetaDTOFromRequest(req requests.CreateMetaRequest) *CreateMetaDTO {
	return &CreateMetaDTO{
		Name:          req.Name,
		PhoneNumberId: req.PhoneNumberId,
		BusinessId:    req.BusinessId,
		Assistants:    req.Assistants,
		UserID:        req.UserID,
	}
}

func (dto *CreateMetaDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}
	if dto.PhoneNumberId == "" {
		return errors.New("O PhoneNumberId é obrigatório e não pode estar em branco.")
	}

	if dto.BusinessId == "" {
		return errors.New("O BusinessId é obrigatório e não pode estar em branco.")
	}
	if dto.UserID == "" {
		return errors.New("O UserID é obrigatório e não pode estar em branco.")
	}

	return nil
}

func (dto *CreateMetaDTO) ToMeta() models.Meta {
	return models.Meta{
		Name:          dto.Name,
		PhoneNumberId: dto.PhoneNumberId,
		BusinessId:    dto.BusinessId,
		UserID:        dto.UserID,
		Assistants:    dto.Assistants,
	}
}
