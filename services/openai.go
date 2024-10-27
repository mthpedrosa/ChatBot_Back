package services

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"

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

func (o *OpenAi) Insert(ctx context.Context, dt *dto.AssistantCreateDTO, userID string) (string, error) {
	assistant := dt.ToAssistant()
	// Log inicial da criação do assistente
	o.logger.Debugf("Create Assistant DTO: %+v", *dt)

	// Verifica se o assistente é do tipo "sub"
	if assistant.Type == "sub" {

		assistant.Subs = []models.Subs{}

		// Apenas insere no MongoDB sem passar pelo OpenAI
		o.logger.Debugf("Creating sub-assistant locally in MongoDB")

		assistant.UserID = userID // Atribui o ID do usuário ao DTO

		// Insere o assistente no MongoDB
		idMongo, err := o.openaiMongo.Insert(ctx, assistant)
		if err != nil {
			return "", err
		}

		o.logger.Infof("Sub-assistant created in MongoDB with ID: %s", idMongo)

		return idMongo, nil
	}

	// Caso contrário, cria no OpenAI primeiro
	o.logger.Debugf("Creating assistant in OpenAI")

	idCriado, err := o.openaiRepository.CreateAssistant(ctx, assistant, "gpt-3.5-turbo-16k")
	if err != nil {
		return "", err
	}

	// Atualiza o DTO com o ID retornado do OpenAI
	assistant.IdAssistant = idCriado.ID // Atribui o ID retornado pelo OpenAI
	assistant.UserID = userID           // Atribui o ID do usuário ao DTO

	// Insere no MongoDB o assistente completo
	idMongo, err := o.openaiMongo.Insert(ctx, assistant)
	if err != nil {
		return "", err
	}

	o.logger.Infof("Assistant created in MongoDB with ID: %s", idMongo)

	return idMongo, nil
}

func (o *OpenAi) Edit(ctx context.Context, dt *dto.AssistantCreateDTO, id string) (string, error) {
	assistant := dt.ToAssistant()

	return o.openaiRepository.UpdateAssistant(ctx, id, "gpt-3.5-turbo-16k", assistant)
}

func (o *OpenAi) FindAll(ctx context.Context, order string, limit int) ([]models.Assistant, error) {
	return o.openaiRepository.ListAssistants(ctx, order, limit)
}

func (o *OpenAi) FindId(ctx context.Context, id string) (models.CreateAssistant, error) {

	// assitant, erro := o.openaiRepository.GetAssistant(ctx, id)
	// if erro != nil {
	// 	return models.Assistant{}, erro
	// }

	assitant, erro := o.openaiMongo.GetAssistant(ctx, id)
	if erro != nil {
		return models.CreateAssistant{}, erro
	}

	return *assitant, nil
}

func (o *OpenAi) Delete(ctx context.Context, id string) (string, error) {

	// Verificamos o tipo do assistante, se for sub sómente deletamos do banco, se for ass deletamos do openAI
	assitant, erro := o.openaiMongo.GetAssistant(ctx, id)
	if erro != nil {
		return "", erro
	}

	erro = o.openaiMongo.Delete(ctx, assitant.ID.Hex())
	if erro != nil {
		return "", erro
	}

	if assitant.Type == "sub" {
		return "Deletado com sucesso", nil
	}

	return o.openaiRepository.DeleteAssistant(ctx, assitant.IdAssistant)
}

// user
func (o *OpenAi) FindAllUser(ctx context.Context, id string) ([]models.CreateAssistant, error) {
	assistants, erro := o.openaiMongo.FindAllUser(ctx, id)
	if erro != nil {
		return nil, erro
	}

	return assistants, nil
}
