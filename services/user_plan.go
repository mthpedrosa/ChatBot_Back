package services

import (
	"autflow_back/models"
	"autflow_back/repositories"
	"autflow_back/requests"
	"context"
	"errors"
	"time"
)

type UserPlanService struct {
	userPlanRepo *repositories.UserPlanRepository
}

func NewUserPlanService(userPlanRepo *repositories.UserPlanRepository) *UserPlanService {
	return &UserPlanService{userPlanRepo: userPlanRepo}
}

// Insert um novo plano de usuário
func (s *UserPlanService) Insert(ctx context.Context, userPlan requests.UserPlanRequest) (string, error) {
	newPlan := models.UserPlan{
		UserID:       userPlan.UserID,
		PlanType:     userPlan.PlanType,
		Subscription: &userPlan.Subscription,
		Credit:       &userPlan.Credit,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	insertedID, err := s.userPlanRepo.Insert(ctx, newPlan)
	if err != nil {
		return "", err
	}

	return insertedID, nil
}

// Edit edits an existing user plan with the provided updateData fields.
func (s *UserPlanService) Edit(ctx context.Context, id string, updateData requests.UserPlanRequest) error {

	newPlan := models.UserPlan{
		UserID:       updateData.UserID,
		PlanType:     updateData.PlanType,
		Subscription: &updateData.Subscription,
		Credit:       &updateData.Credit}

	// Atualiza o plano no repositório com os campos relevantes
	return s.userPlanRepo.Edit(ctx, id, newPlan)
}

func (s *UserPlanService) FindId(ctx context.Context, id string) (models.UserPlan, error) {
	userPlan, erro := s.userPlanRepo.FindId(ctx, id)
	if erro != nil {
		return models.UserPlan{}, nil

	}

	return *userPlan, nil
}

// Remove um plano de usuário
func (s *UserPlanService) Delete(ctx context.Context, id string) error {
	// Chama a função de deletar no repositório
	err := s.userPlanRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Decrementa o saldo de mensagens restantes em um plano de assinatura
func (s *UserPlanService) DecrementMessagesRemaining(ctx context.Context, id string, messagesToDecrement int) error {
	userPlan, err := s.userPlanRepo.FindId(ctx, id)
	if err != nil {
		return err
	}

	// Verifica se o plano é do tipo assinatura e se há mensagens restantes suficientes
	if userPlan.PlanType == "subscription" && userPlan.Subscription != nil {
		if userPlan.Subscription.MessagesRemaining < messagesToDecrement {
			return errors.New("mensagens insuficientes no plano de assinatura")
		}
		userPlan.Subscription.MessagesRemaining -= messagesToDecrement
		userPlan.UpdatedAt = time.Now()

		return s.userPlanRepo.Edit(ctx, id, *userPlan)
	}

	return errors.New("plano não é do tipo assinatura ou informações de assinatura ausentes")
}

// Decrementa o saldo de créditos em um plano de crédito
func (s *UserPlanService) DecrementCreditBalance(ctx context.Context, id string, messagesCount int) error {
	userPlan, err := s.userPlanRepo.FindId(ctx, id)
	if err != nil {
		return err
	}

	// Verifica se o plano é do tipo crédito e se o saldo é suficiente
	if userPlan.PlanType == "credit" && userPlan.Credit != nil {
		totalCost := float64(messagesCount) * userPlan.Credit.CostPerMessage
		if userPlan.Credit.Balance < totalCost {
			return errors.New("saldo insuficiente no plano de crédito")
		}
		userPlan.Credit.Balance -= totalCost
		userPlan.UpdatedAt = time.Now()

		return s.userPlanRepo.Edit(ctx, id, *userPlan)
	}

	return errors.New("plano não é do tipo crédito ou informações de crédito ausentes")
}
