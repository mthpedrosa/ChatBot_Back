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
	metaRepository   *repositories.Metas
}

func NewOpenAi(openai interfaces.OpenAIClientRepository, openaiMongo *repositories.OpenaiMongo,
	logger utils.Logger, metaRepository *repositories.Metas) *OpenAi {
	return &OpenAi{
		logger:           logger,
		openaiRepository: openai,
		openaiMongo:      openaiMongo,
		metaRepository:   metaRepository,
	}
}

func (o *OpenAi) Insert(ctx context.Context, dt *dto.AssistantCreateDTO, userID string) (string, error) {
	assistant := dt.ToAssistant()
	// Log inicial da criação do assistente
	fmt.Printf("Create Assistant DTO: %+v", *dt)

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

	// Verifica se o assistente sendo editado está sendo ativado e é do tipo "ass"
	if assistant.Type == "ass" && assistant.Active {
		// Desativa qualquer outro assistente "ass" ativo do mesmo usuário
		err := o.openaiMongo.DeactivateOtherAssistants(ctx, "", "ass", userID)
		if err != nil {
			return "", err
		}
	}

	// Caso contrário, cria no OpenAI primeiro
	fmt.Printf("Creating assistant in OpenAI")

	assistant.Instructions += assistant.Info

	//Verificamos se existe subs vinculado
	if len(dt.Subs) > 0 {
		for i, sub := range dt.Subs {
			fmt.Printf("Subs[%d] - MongoID: %s\n", i, sub.MongoID)
			subGet, error := o.openaiMongo.GetAssistant(ctx, sub.MongoID)
			if error != nil {
				return "", error
			}

			assistant.Instructions += " topicos: " + subGet.Instructions + "info" + subGet.Info
		}
	}

	idCriado, err := o.openaiRepository.CreateAssistant(ctx, assistant, "gpt-3.5-turbo-16k")
	if err != nil {
		return "", err
	}

	assistant.IdAssistant = idCriado.ID
	assistant.UserID = userID

	idMongo, err := o.openaiMongo.Insert(ctx, assistant)
	if err != nil {
		return "", err
	}

	fmt.Printf("Assistant created in MongoDB with ID: %s", idMongo)

	return idMongo, nil
}

func (o *OpenAi) Edit(ctx context.Context, dt *dto.AssistantCreateDTO, id string) (string, error) {
	assistant := dt.ToAssistant()

	fmt.Print("---------------------")
	fmt.Println(assistant)

	// Verifica se o assistente sendo editado está sendo ativado e é do tipo "ass"
	if assistant.Type == "ass" && assistant.Active {
		// Desativa qualquer outro assistente "ass" ativo do mesmo usuário
		err := o.openaiMongo.DeactivateOtherAssistants(ctx, id, "ass", assistant.UserID)
		if err != nil {
			return "", err
		}
	}

	// Insere a atualização no MongoDB
	err := o.openaiMongo.Edit(ctx, id, assistant)
	if err != nil {
		return "", err
	}

	newAssistante, err := o.openaiMongo.GetAssistant(ctx, id)
	fmt.Println(newAssistante)

	// Verifica se o assistente é do tipo "ass"
	if newAssistante.Type == "ass" {
		fmt.Printf("Updating assistant in OpenAI")
		newAssistante.Instructions += newAssistante.Info

		// Verifica se existem subs vinculados
		if len(dt.Subs) > 0 {
			for i, sub := range dt.Subs {
				fmt.Printf("Subs[%d] - MongoID: %s\n", i, sub.MongoID)
				subGet, err := o.openaiMongo.GetAssistant(ctx, sub.MongoID)
				if err != nil {
					return "", err
				}

				newAssistante.Instructions += "topicos: " + subGet.Instructions + "info" + subGet.Info
			}
		}

		// Atualiza no OpenAI antes de atualizar no MongoDB
		updatedID, err := o.openaiRepository.UpdateAssistant(ctx, newAssistante.IdAssistant, "gpt-3.5-turbo-16k", *newAssistante)
		if err != nil {
			return "", err
		}

		// Atualiza o ID do assistente para refletir no MongoDB após edição no OpenAI
		assistant.IdAssistant = updatedID
	}

	return "Editado com sucesso", nil
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

func (o *OpenAi) UpdateMeta(ctx context.Context, id string) (string, error) {
	assistants, err := o.openaiMongo.FindAllUser(ctx, id)
	if err != nil {
		return "", err
	}

	contaMeta, err := o.metaRepository.Find(ctx, "user_id="+id)
	if err != nil {
		return "", err
	}

	newAssistantes := []models.AssistantIds{}
	for _, assistant := range assistants {
		newAssistantes = append(newAssistantes, models.AssistantIds{
			OpenId: assistant.IdAssistant,
			Id:     assistant.ID.Hex(),
			Active: assistant.Active,
		})
	}

	contaMeta[0].Assistants = newAssistantes

	// Atualiza a conta meta com a nova lista de assistentes
	err = o.metaRepository.Edit(ctx, contaMeta[0].ID.Hex(), contaMeta[0])
	if err != nil {
		return "", err
	}

	return "Conta meta atualizada com sucesso", nil
}
