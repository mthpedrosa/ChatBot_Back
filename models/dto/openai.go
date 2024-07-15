package dto

import "autflow_back/models"

type CreateAssistantDTO struct {
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Instructions string `json:"instructions,omitempty" bson:"instructions,omitempty"`
}

func (dto *CreateAssistantDTO) ToMeta() models.Assistant {
	return models.Assistant{
		Name:         dto.Name,
		Instructions: dto.Instructions,
	}
}
