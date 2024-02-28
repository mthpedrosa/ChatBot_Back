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

type Meta struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty" bson:"name,omitempty"`
	Token          string             `json:"token" bson:"token"`
	MetaID         string             `json:"meta_id" bson:"meta_id"`
	PhonesMeta     []PhonesMeta       `json:"phones_meta" bson:"phones_meta"`
	EditPermission []string           `json:"edit_permission" bson:"edit_permission"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt       time.Time          `json:"update_at" bson:"update_at"`
	CreatedBy      string             `json:"created_by_id" bson:"created_by_id"`
	Webhook        string             `json:"webhook" bson:"webhook"`
	OtherFields    []Fields           `json:"other_fields" bson:"other_fields"`
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
	if meta.Token == "" {
		return errors.New("O token é obrigatório e não pode estar em branco.")
	}

	if meta.MetaID == "" {
		return errors.New("O ID da conta META é obrigatório e não pode estar em branco.")
	}

	for _, telefones := range meta.PhonesMeta {
		if telefones.Id == "" {
			return errors.New("O id do telefone é obrigatorio")
		}
	}

	return nil
}

// Formatar - Remove white spaces.
func (meta *Meta) format(etapa string) error {
	meta.Name = strings.TrimSpace(meta.Name)
	meta.Token = strings.TrimSpace(meta.Token)

	/*if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}

		usuario.Senha = string(senhaComHash)
	}*/

	return nil
}
