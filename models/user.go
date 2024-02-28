package models

import (
	"autflow_back/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	Password   string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt   time.Time          `json:"update_at" bson:"update_at"`
	LastActive time.Time          `json:"last_active,omitempty" bson:"last_active"`
	Profile    string             `json:"profile,omitempty" bson:"profile,omitempty"`
}

func (user *User) Prepare(etapa string) error {
	if erro := user.validate(etapa); erro != nil {
		return erro
	}

	if erro := user.format(etapa); erro != nil {
		return erro
	}

	return nil
}

func (user *User) validate(etapa string) error {
	if user.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}
	if user.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco.")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("E-mail invalido")
	}

	if etapa == "cadastro" && user.Password == "" {
		return errors.New("A senha é obrigatório e não pode estar em branco.")
	}
	if user.Profile == "" {
		return errors.New("O perfil é obrigatório e não pode estar em branco.")
	}

	return nil
}

func (user *User) format(etapa string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if etapa == "cadastro" {
		senhaComHash, erro := security.Hash(user.Password)
		if erro != nil {
			return erro
		}

		user.Password = string(senhaComHash)
	}

	return nil
}
