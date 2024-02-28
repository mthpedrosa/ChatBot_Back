package dto

import (
	"autflow_back/models"
	"errors"
	"time"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCustomerDTO struct {
	Name        string          `json:"name,omitempty" bson:"name,omitempty"`
	Email       string          `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string          `json:"phone,omitempty" bson:"phone,omitempty"`
	WhatsAppID  string          `json:"whatsapp_id,omitempty" bson:"whatsapp_id,omitempty"`
	OtherFields []models.Fields `json:"other_fields" bson:"other_fields"`
}

type CustomerListDTO struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	WhatsAppID  string             `json:"whatsapp_id,omitempty" bson:"whatsapp_id,omitempty"`
	OtherFields []models.Fields    `json:"other_fields" bson:"other_fields"`
}

type CustomerDetailDTO struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	WhatsAppID  string             `json:"whatsapp_id,omitempty" bson:"whatsapp_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	OtherFields []models.Fields    `json:"other_fields" bson:"other_fields"`
}

func (dto *CreateCustomerDTO) ToCustomer() models.Customer {
	return models.Customer{
		Name:        dto.Name,
		Email:       dto.Email,
		Phone:       dto.Phone,
		WhatsAppID:  dto.WhatsAppID,
		OtherFields: dto.OtherFields,
	}
}

func (dto *CreateCustomerDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}
	if dto.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco.")
	}

	if dto.WhatsAppID == "" {
		return errors.New("O WA_ID é obrigatório e não pode estar em branco.")
	}

	if erro := checkmail.ValidateFormat(dto.Email); erro != nil {
		return errors.New("E-mail invalido")
	}

	return nil
}
