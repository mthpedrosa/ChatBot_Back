package dto

import (
	"autflow_back/models"
	"autflow_back/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserDTO struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Profile  string `json:"profile,omitempty" bson:"profile,omitempty"`
}

type UserListDTO struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Email   string             `json:"email,omitempty" bson:"email,omitempty"`
	Profile string             `json:"profile,omitempty" bson:"profile,omitempty"`
}

type UserDetailDTO struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	Password   string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt   time.Time          `json:"update_at" bson:"update_at"`
	LastActive time.Time          `json:"last_active,omitempty" bson:"last_active,omitempty"`
	Profile    string             `json:"profile,omitempty" bson:"profile,omitempty"`
}

func (dto *CreateUserDTO) ToUser() models.User {
	return models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Profile:  dto.Profile,
	}
}

func (user *CreateUserDTO) Prepare(etapa string) error {
	if erro := user.validate(etapa); erro != nil {
		return erro
	}

	if erro := user.format(etapa); erro != nil {
		return erro
	}

	return nil
}

func (user *CreateUserDTO) validate(etapa string) error {
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

func (user *CreateUserDTO) format(etapa string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if etapa == "cadastro" {
		passwordHash, erro := security.Hash(user.Password)
		if erro != nil {
			return erro
		}

		user.Password = string(passwordHash)
	}

	return nil
}
