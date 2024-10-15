package services

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"fmt"

	"autflow_back/utils"
	"context"
)

type OpenAi struct {
	openaiRepository interfaces.OpenAIClientRepository
	openaiMongo      *repositories.OpenaiMongo
	logger           utils.Logger
}

func NewOpenAi(openai interfaces.OpenAIClientRepository, openaiMongo *repositories.OpenaiMongo,
	logger utils.Logger) *OpenAi {
	return &OpenAi{
		logger:           logger,
		openaiRepository: openai,
		openaiMongo:      openaiMongo,
	}
}

func (o *OpenAi) Insert(ctx context.Context, dt *dto.CreateAssistantDTO, userID string) (models.Assistant, error) {
	assistant := *dt

	o.logger.Debugf("Create Assistant: %+v", assistant)

	idCriado, erro := o.openaiRepository.CreateAssistant(ctx, assistant, "gpt-3.5-turbo-16k")
	if erro != nil {
		return models.Assistant{}, erro
	}

	// add create in local db
	idCriado.UserID = userID
	idMongo, erro := o.openaiMongo.Insert(ctx, *idCriado)
	if erro != nil {
		return models.Assistant{}, erro
	}

	fmt.Println("ID MONGO ASSISTANT:", idMongo)

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

// user
func (o *OpenAi) FindAllUser(ctx context.Context, id string) ([]models.Assistant, error) {
	assistants, erro := o.openaiMongo.FindAllUser(ctx, id)
	if erro != nil {
		return nil, erro
	}

	return assistants, nil
}
