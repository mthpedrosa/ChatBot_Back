package services

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"autflow_back/utils"
	"context"
)

type Conversations struct {
	conversarionsRepository *repositories.Conversations
	logger                  utils.Logger
}

func NewConversation(
	conversarionsRepository *repositories.Conversations,
	logger utils.Logger) *Conversations {
	return &Conversations{
		logger:                  logger,
		conversarionsRepository: conversarionsRepository,
	}
}

func (r *Conversations) Insert(ctx context.Context, dt *dto.ConversationCreateDTO) (string, error) {
	conversation := dt.ToConversation()

	idCriado, erro := r.conversarionsRepository.Insert(ctx, conversation)
	if erro != nil {
		return "", erro
	}

	return idCriado, nil
}

func (r *Conversations) Find(ctx context.Context, query string) ([]models.Conversation, error) {
	conversations, erro := r.conversarionsRepository.Find(ctx, query)
	if erro != nil {
		return nil, erro
	}

	return conversations, nil
}

func (r *Conversations) FindId(ctx context.Context, id string) (models.Conversation, error) {
	conversation, erro := r.conversarionsRepository.FindId(ctx, id)
	if erro != nil {
		return models.Conversation{}, erro
	}

	return *conversation, nil
}

func (r *Conversations) Edit(ctx context.Context, id string, dt *dto.ConversationCreateDTO) error {
	conversation := dt.ToConversation()
	return r.conversarionsRepository.Edit(ctx, id, conversation)
}

func (r *Conversations) Delete(ctx context.Context, id string) error {
	return r.conversarionsRepository.Delete(ctx, id)
}
