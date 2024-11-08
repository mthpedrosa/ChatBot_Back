package services

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/repositories"
	"autflow_back/requests"
	"autflow_back/utils"
	"context"
	"fmt"
)

type Reports struct {
	metaRepository          *repositories.Metas
	customerRepository      *repositories.Customers
	sessionRepository       *repositories.Session
	conversarionsRepository *repositories.Conversations
	logger                  utils.Logger
	openaiRepository        interfaces.OpenAIClientRepository
	whatsappRepository      interfaces.WhatsappRepository
}

func NewReports(
	meta *repositories.Metas,
	customer *repositories.Customers,
	session *repositories.Session,
	conversarions *repositories.Conversations,
	logger utils.Logger, openai interfaces.OpenAIClientRepository, whatsapp interfaces.WhatsappRepository) *Reports {
	return &Reports{
		metaRepository:          meta,
		logger:                  logger,
		customerRepository:      customer,
		sessionRepository:       session,
		conversarionsRepository: conversarions,
		openaiRepository:        openai,
		whatsappRepository:      whatsapp,
	}
}

func (o *Reports) Cost(ctx context.Context, dt *requests.CostParams) ([]map[string]interface{}, error) {
	fmt.Print("dentro de service")

	// Converte as datas de strings para time.Time
	startDate, endDate, err := dt.ParseDates()
	if err != nil {
		fmt.Print("erro ao fazer parse")
		return nil, err
	}

	fmt.Print("PASSA DATA de service")

	// Consultamos a conta meta  do id mandado
	meta, err := o.metaRepository.FindId(ctx, dt.MetaId)
	if err != nil {
		return nil, err
	}

	// Array para armazenar as informações detalhadas de cada assistente
	assistantsReport := []map[string]interface{}{}

	// Itera sobre os assistentes vinculados à conta Meta
	for _, assistant := range meta.Assistants {
		fmt.Printf("Consultando custo do assistente: %v\n", assistant)

		// Chamar o método `FindCost` do repositório de sessões
		sessions, err := o.sessionRepository.FindCost(ctx, startDate, endDate, assistant.Id)
		if err != nil {
			return nil, fmt.Errorf("erro ao consultar sessões: %v", err)
		}

		// Inicializa o custo e contadores
		totalCost := 0.0
		totalMessages := 0

		// Processa cada sessão para calcular o custo e contar mensagens
		for _, session := range sessions {
			sessionCost, messages, err := o.calculateSessionCost(ctx, session)
			if err != nil {
				return nil, err
			}
			totalCost += sessionCost
			totalMessages += messages
		}

		// Cria um objeto com o resumo das informações para o assistente
		assistantReport := map[string]interface{}{
			"assistant_id":   assistant.Id,
			"sessions_count": len(sessions),
			"messages_count": totalMessages,
			"total_cost":     totalCost,
		}

		// Adiciona o relatório do assistente ao array principal
		assistantsReport = append(assistantsReport, assistantReport)

		fmt.Printf("Relatório do assistente: %v\n", assistantReport)
	}

	return assistantsReport, nil
}

func (o *Reports) calculateSessionCost(ctx context.Context, session models.Session) (float64, int, error) {
	conversation, err := o.conversarionsRepository.FindId(ctx, session.ConversationId)
	if err != nil {
		return 0, 0, err
	}

	costPerMessage := float64(len(conversation.Messages)) * 0.01
	totalCost := costPerMessage + 0.15

	return totalCost, len(conversation.Messages), nil
}
