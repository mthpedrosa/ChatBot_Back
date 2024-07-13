package services

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/utils"
	"context"
)

type OpenAi struct {
	openaiRepository interfaces.OpenAIClientRepository
	logger           utils.Logger
}

func NewOpenAi(openai interfaces.OpenAIClientRepository,
	logger utils.Logger) *OpenAi {
	return &OpenAi{
		logger:           logger,
		openaiRepository: openai,
	}
}

func (o *OpenAi) Insert(ctx context.Context, dt *dto.CreateAssistantDTO, idClient string) (models.Assistant, error) {
	assistant := *dt

	o.logger.Debugf("Create Assistant: %+v", assistant)

	idCriado, erro := o.openaiRepository.CreateAssistant(ctx, assistant, "gpt-3.5-turbo-16k")
	if erro != nil {
		return models.Assistant{}, erro
	}

	return *idCriado, nil
}
