package services

import (
	"autflow_back/models"
	"autflow_back/repositories"
	"autflow_back/requests"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPlanService struct {
	userPlanRepo *repositories.UserPlanRepository
}

func NewUserPlanService(userPlanRepo *repositories.UserPlanRepository) *UserPlanService {
	return &UserPlanService{userPlanRepo: userPlanRepo}
}

// Insert um novo plano de usuário
func (s *UserPlanService) Insert(ctx context.Context, userPlan requests.UserPlanRequest) (primitive.ObjectID, error) {
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
		return primitive.NilObjectID, err
	}

	return insertedID, nil
}

// Edit edits an existing user plan with the provided updateData fields.
func (s *UserPlanService) Edit(ctx context.Context, planID primitive.ObjectID, updateData models.UserPlan) error {
	// Verifica se o plano existe
	existingPlan, err := s.userPlanRepo.FindByID(ctx, planID)
	if err != nil {
		return err
	}

	// Cria um mapa de atualizações para os campos modificados
	updates := bson.M{}

	if updateData.PlanType != "" && updateData.PlanType != existingPlan.PlanType {
		updates["plan_type"] = updateData.PlanType
	}
	if updateData.Subscription != nil {
		updates["subscription"] = updateData.Subscription
	}
	if updateData.Credit != nil {
		updates["credit"] = updateData.Credit
	}
	if len(updates) > 0 { // Caso haja algo a atualizar
		updates["updated_at"] = time.Now() // Atualiza o campo updated_at
	}

	// Verifica se há atualizações antes de chamar o repositório
	if len(updates) == 0 {
		return errors.New("nenhuma alteração detectada para o plano do usuário")
	}

	// Atualiza o plano no repositório com os campos relevantes
	return s.userPlanRepo.Edit(ctx, planID, updates)
}

func (s *UserPlanService) FindId(ctx context.Context, id string) (models.UserPlan, error) {
	userPlan, erro := s.userPlanRepo.FindId(ctx, id)
	if erro != nil {
		return models.UserPlan{}, nil

	}

	return *userPlan, nil
}

// Remove um plano de usuário
func (s *UserPlanService) Delete(ctx context.Context, planID primitive.ObjectID) error {
	// Chama a função de deletar no repositório
	err := s.userPlanRepo.Delete(ctx, planID)
	if err != nil {
		return err
	}

	return nil
}

// Decrementa o saldo de mensagens restantes em um plano de assinatura
func (s *UserPlanService) DecrementMessagesRemaining(ctx context.Context, planID primitive.ObjectID, messagesToDecrement int) error {
	userPlan, err := s.userPlanRepo.FindByID(ctx, planID)
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

		// Cria um mapa de atualizações para os campos modificados
		updates := bson.M{}
		updates["subscription"] = userPlan.Subscription
		updates["updated_at"] = time.Now()

		return s.userPlanRepo.Edit(ctx, userPlan.ID, updates)
	}

	return errors.New("plano não é do tipo assinatura ou informações de assinatura ausentes")
}

// Decrementa o saldo de créditos em um plano de crédito
func (s *UserPlanService) DecrementCreditBalance(ctx context.Context, planID primitive.ObjectID, messagesCount int) error {
	userPlan, err := s.userPlanRepo.FindByID(ctx, planID)
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

		// Cria um mapa de atualizações para os campos modificados
		updates := bson.M{}
		updates["credit"] = userPlan.Credit
		updates["updated_at"] = time.Now()

		return s.userPlanRepo.Edit(ctx, userPlan.ID, updates)
	}

	return errors.New("plano não é do tipo crédito ou informações de crédito ausentes")
}
