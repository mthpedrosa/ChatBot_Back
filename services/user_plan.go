package services

import (
	"autflow_back/models"
	"autflow_back/repositories"
	"autflow_back/requests"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
		Subscription: userPlan.Subscription,
		Credit:       userPlan.Credit,
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

	updateFields := bson.M{}

	if updateData.UserID != "" {
		updateFields["user_id"] = updateData.UserID
	}
	if updateData.PlanType != "" {
		updateFields["plan_type"] = updateData.PlanType
	}

	// Verifica se Subscription deve ser atualizado
	if updateData.Subscription != (models.SubscriptionPlan{}) {
		updateFields["subscription"] = updateData.Subscription
	}

	// Verifica se Credit deve ser atualizado
	if updateData.Credit != (models.CreditPlan{}) {
		updateFields["credit"] = updateData.Credit
	}
	// Atualiza o plano no repositório com os campos relevantes
	return s.userPlanRepo.Edit(ctx, id, updateFields)
}

func (s *UserPlanService) Find(ctx context.Context, query string) ([]models.UserPlan, error) {
	usersPlan, erro := s.userPlanRepo.Find(ctx, query)
	if erro != nil {
		return nil, erro
	}

	return usersPlan, nil
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

// DecrementMessagesRemaining decrementa o saldo de mensagens restantes em um plano de assinatura
func (s *UserPlanService) DecrementMessagesRemaining(ctx context.Context, id string, messagesToDecrement int) error {
	// Busca o plano do usuário pelo ID
	userPlan, err := s.userPlanRepo.FindId(ctx, id)
	if err != nil {
		fmt.Println("Erro ao buscar plano do usuário:", err) // Debug
		return err
	}

	fmt.Println("Verificando tipo de plano e saldo de mensagens...") // Debug

	// Verifica se o plano é do tipo assinatura e possui saldo suficiente de mensagens
	if userPlan.PlanType == "subscription" && userPlan.Subscription != (models.SubscriptionPlan{}) {
		if userPlan.Subscription.MessagesRemaining < messagesToDecrement {
			fmt.Println("Saldo insuficiente de mensagens no plano de assinatura") // Debug
			return errors.New("mensagens insuficientes no plano de assinatura")
		}

		// Calcula o novo saldo de mensagens
		newMessagesRemaining := userPlan.Subscription.MessagesRemaining - messagesToDecrement

		// Prepara o update específico
		update := bson.M{
			"subscription.messages_remaining": newMessagesRemaining,
		}

		fmt.Println("Dentro da onde tira a mensagem - preparando para atualizar") // Debug

		// Executa a atualização no repositório
		err := s.userPlanRepo.Edit(ctx, id, update)
		if err != nil {
			fmt.Println("Erro ao atualizar mensagens restantes no plano:", err) // Debug
			return fmt.Errorf("erro ao atualizar mensagens restantes no plano: %v", err)
		}

		fmt.Println("Saldo de mensagens atualizado com sucesso!") // Debug
		return nil
	}

	fmt.Println("Plano não é do tipo assinatura ou informações de assinatura ausentes") // Debug
	return errors.New("plano não é do tipo assinatura ou informações de assinatura ausentes")
}

// DecrementCreditBalance decrementa o saldo de créditos em um plano de crédito
func (s *UserPlanService) DecrementCreditBalance(ctx context.Context, id string, messagesCount int) error {
	userPlan, err := s.userPlanRepo.FindId(ctx, id)
	if err != nil {
		return err
	}

	// Verifica se o plano é do tipo crédito e se o saldo é suficiente
	if userPlan.PlanType == "credit" && userPlan.Credit != (models.CreditPlan{}) {
		totalCost := float64(messagesCount) * userPlan.Credit.CostPerMessage
		if userPlan.Credit.Balance < totalCost {
			return errors.New("saldo insuficiente no plano de crédito")
		}

		// Calcula o novo saldo
		newBalance := userPlan.Credit.Balance - totalCost

		// Cria o objeto de atualização
		update := bson.M{
			"credit.balance": newBalance,
		}

		// Atualiza o campo específico no repositório
		return s.userPlanRepo.Edit(ctx, id, update)
	}

	return errors.New("plano não é do tipo crédito ou informações de crédito ausentes")
}
