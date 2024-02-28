package dto

import (
	"autflow_back/models"
	"autflow_back/requests"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateMetaDTO struct {
	Name           string              `json:"name,omitempty" bson:"name,omitempty"`
	Token          string              `json:"token" bson:"token"`
	MetaID         string              `json:"meta_id" bson:"meta_id"`
	PhonesMeta     []models.PhonesMeta `json:"phones_meta" bson:"phones_meta"`
	EditPermission []string            `json:"edit_permission" bson:"edit_permission"`
	CreatedBy      string              `json:"created_by_id" bson:"created_by_id"`
}

type MetaListDTO struct {
	ID             primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string              `json:"name,omitempty" bson:"name,omitempty"`
	MetaID         string              `json:"meta_id" bson:"meta_id"`
	PhonesMeta     []models.PhonesMeta `json:"phones_meta" bson:"phones_meta"`
	EditPermission []string            `json:"edit_permission" bson:"edit_permission"`
}

type MetaDetailDTO struct {
	ID             primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string              `json:"name,omitempty" bson:"name,omitempty"`
	Token          string              `json:"token" bson:"token"`
	MetaID         string              `json:"meta_id" bson:"meta_id"`
	PhonesMeta     []models.PhonesMeta `json:"phones_meta" bson:"phones_meta"`
	EditPermission []string            `json:"edit_permission" bson:"edit_permission"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdateAt       time.Time           `json:"update_at" bson:"update_at"`
	CreatedBy      string              `json:"created_by_id" bson:"created_by_id"`
	Webhook        string              `json:"webhook" bson:"webhook"`
	OtherFields    []models.Fields     `json:"other_fields" bson:"other_fields"`
}

type PhonesMeta struct {
	Id     string `json:"id" bson:"id"`
	Number string `json:"number" bson:"number"`
	Name   string `json:"name" bson:"name"`
}

func NewCreateMetaDTOFromRequest(req requests.CreateMetaRequest) *CreateMetaDTO {
	phonesMeta := make([]PhonesMeta, len(req.PhonesMeta))
	for i, phone := range req.PhonesMeta {
		phonesMeta[i] = PhonesMeta{
			Id:     phone.Id,
			Number: phone.Number,
			Name:   phone.Name,
		}
	}

	return &CreateMetaDTO{
		Name:           req.Name,
		Token:          req.Token,
		MetaID:         req.MetaID,
		PhonesMeta:     req.PhonesMeta,
		EditPermission: req.EditPermission,
	}
}

func (dto *CreateMetaDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}
	if dto.Token == "" {
		return errors.New("O token é obrigatório e não pode estar em branco.")
	}

	if dto.MetaID == "" {
		return errors.New("O ID da conta META é obrigatório e não pode estar em branco.")
	}

	for _, phone := range dto.PhonesMeta {
		if phone.Id == "" {
			return errors.New("O id do telefone é obrigatório.")
		}
	}

	return nil
}

func (dto *CreateMetaDTO) ToMeta() models.Meta {
	return models.Meta{
		Name:           dto.Name,
		Token:          dto.Token,
		MetaID:         dto.MetaID,
		PhonesMeta:     dto.PhonesMeta,
		EditPermission: dto.EditPermission,
	}
}
