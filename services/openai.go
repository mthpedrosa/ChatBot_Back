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

func (o *OpenAi) Edit(ctx context.Context, dt *dto.CreateAssistantDTO, id string) (string, error) {
	assistant := dt.ToMeta()

	return o.openaiRepository.UpdateAssistant(ctx, id, "gpt-3.5-turbo-16k", assistant)
}

func (o *OpenAi) FindAll(ctx context.Context, order string, limit int) ([]models.Assistant, error) {
	return o.openaiRepository.ListAssistants(ctx, order, limit)
}

func (o *OpenAi) FindId(ctx context.Context, id string) (models.Assistant, error) {

	assitant, erro := o.openaiRepository.GetAssistant(ctx, id)
	if erro != nil {
		return models.Assistant{}, erro
	}

	return *assitant, nil
}

func (o *OpenAi) Delete(ctx context.Context, id string) (string, error) {
	return o.openaiRepository.DeleteAssistant(ctx, id)
}
