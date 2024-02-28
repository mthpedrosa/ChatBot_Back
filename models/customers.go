package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	WhatsAppID  string             `json:"whatsapp_id,omitempty" bson:"whatsapp_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	OtherFields []Fields           `json:"other_fields" bson:"other_fields"`
}

// Prepare - Call user validation and formatting methods.
func (customer *Customer) Prepare() error {
	if erro := customer.validate(); erro != nil {
		return erro
	}

	if erro := customer.format(); erro != nil {
		return erro
	}

	return nil
}

// Validate - Check if the fields are empty.
func (customer *Customer) validate() error {
	if customer.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}
	if customer.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco.")
	}

	if customer.WhatsAppID == "" {
		return errors.New("O WA_ID é obrigatório e não pode estar em branco.")
	}

	if erro := checkmail.ValidateFormat(customer.Email); erro != nil {
		return errors.New("E-mail invalido")
	}

	return nil
}

// Format - Remove white spaces.
func (customer *Customer) format() error {
	customer.Name = strings.TrimSpace(customer.Name)
	customer.Email = strings.TrimSpace(customer.Email)

	return nil
}
