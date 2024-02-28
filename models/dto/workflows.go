package dto

import (
	"autflow_back/models"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateWorkflowDTO struct {
	PhoneMetaId string        `json:"phone_meta_id" bson:"phone_meta_id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Nodes       []models.Node `json:"nodes" bson:"nodes"`
	Active      string        `json:"active" bson:"active"`
	FirstNode   string        `json:"first_node" bson:"first_node"`
	LastNode    string        `json:"last_node" bson:"last_node"`
	CreatedBy   string        `json:"created_by_id" bson:"created_by_id"`
}

type WorkflowListDTO struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PhoneMetaId string             `json:"phone_meta_id" bson:"phone_meta_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Active      string             `json:"active" bson:"active"`
}

type WorkflowDetailDTO struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PhoneMetaId string             `json:"phone_meta_id" bson:"phone_meta_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Nodes       []models.Node      `json:"nodes" bson:"nodes"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	CreatedBy   string             `json:"created_by_id" bson:"created_by_id"`
	Active      string             `json:"active" bson:"active"`
	FirstNode   string             `json:"first_node" bson:"first_node"`
	LastNode    string             `json:"last_node" bson:"last_node"`
}

type FlowData struct {
	Ctx        context.Context
	Customer   models.Customer
	Session    models.Session
	Message    models.MessagePayload
	MetaTokens models.MetaIds
	Arguments  string
}

func (dto *CreateWorkflowDTO) ToWorkflow() models.Workflow {
	return models.Workflow{
		PhoneMetaId: dto.PhoneMetaId,
		Name:        dto.Name,
		Description: dto.Description,
		Nodes:       dto.Nodes,
		Active:      dto.Active,
		FirstNode:   dto.FirstNode,
		LastNode:    dto.LastNode,
		CreatedBy:   dto.CreatedBy,
	}
}

func (workflow *CreateWorkflowDTO) Prepare(etapa string) error {
	if erro := workflow.validate(etapa); erro != nil {
		return erro
	}

	if erro := workflow.format(etapa); erro != nil {
		return erro
	}

	return nil
}

func (workflow *CreateWorkflowDTO) validate(etapa string) error {
	if workflow.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco.")
	}
	if workflow.PhoneMetaId == "" {
		return errors.New("O Telefone é obrigatório e não pode estar em branco.")
	}

	return nil
}

func (workflow *CreateWorkflowDTO) format(etapa string) error {
	workflow.Name = strings.TrimSpace(workflow.Name)
	workflow.PhoneMetaId = strings.TrimSpace(workflow.PhoneMetaId)

	return nil
}
