package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhonesMeta struct {
	Id         string `json:"id"`
	Number     string `json:"number" bson:"number"`
	Name       string `json:"name" bson:"name"`
	BusinessId string `json:"business_id" bson:"business_id"`
}

type AssistantIds struct {
	OpenId string `json:"open_id"`
	Id     string `json:"id"`
	Active bool   `json:"active"`
}

type Meta struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumberId string             `json:"phone_id" bson:"phone_id"`
	BusinessId    string             `json:"business_id" bson:"business_id"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt      time.Time          `json:"update_at" bson:"update_at"`
	Assistants    []AssistantIds     `json:"assistants" bson:"assistants"`
	UserID        string             `json:"user_id" bson:"user_id"`
	//MetaID     string             `json:"meta_id" bson:"meta_id"`
	//Webhook        string         `json:"webhook" bson:"webhook"`
	//Token          string             `json:"token" bson:"token"`
	//EditPermission []string       `json:"edit_permission" bson:"edit_permission"`
	//PhonesMeta     []PhonesMeta   `json:"phones_meta" bson:"phones_meta"`
}

type MetaIds struct {
	PhoneID    string `json:"phone_id,omitempty" bson:"phone_id,omitempty"`
	Token      string `json:"token,omitempty" bson:"token,omitempty"`
	BusinessId string `json:"business_id" bson:"business_id"`
}

// Prepare - Call user validation and formatting methods.
func (meta *Meta) Prepare(etapa string) error {
	if erro := meta.validate(etapa); erro != nil {
		return erro
	}

	if erro := meta.format(etapa); erro != nil {
		return erro
	}

	return nil
}

// Validate - Check if the fields are empty.
func (meta *Meta) validate(etapa string) error {
	if meta.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}

	if meta.PhoneNumberId == "" {
		return errors.New("O numero é obrigatório e não pode estar em branco.")
	}

	if meta.BusinessId == "" {
		return errors.New("O BusinessId é obrigatório e não pode estar em branco.")
	}

	return nil
}

// Formatar - Remove white spaces.
func (meta *Meta) format(etapa string) error {
	meta.Name = strings.TrimSpace(meta.Name)
	return nil
}
