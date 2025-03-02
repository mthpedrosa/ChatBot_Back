package controllers

import (
	"autflow_back/models"
	"autflow_back/requests"
	"autflow_back/services"
	"autflow_back/src/responses"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

/*type MetaIds struct {
	IdTelefone   string
	TokenConexao string
}*/

type Webhook struct {
	messageHandler *services.MessageHandler
	metaService    *services.Meta
	userPlan       *services.UserPlanService
}

func NewWebhookController(messageHandler *services.MessageHandler, metaService *services.Meta, userPlan *services.UserPlanService) *Webhook {
	return &Webhook{
		messageHandler: messageHandler,
		metaService:    metaService,
		userPlan:       userPlan,
	}
}

func (r *Webhook) WebhookRun(c echo.Context) error {
	var payload models.WebhookPayload
	fmt.Println("Start conversation process")
	fmt.Println(c.Request())

	// Extract the ID from the request
	webhookId := c.Param("id")
	if webhookId == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// Decode the request body (JSON payload) into 'payload'
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return responses.Erro(c, http.StatusBadRequest, errors.New("Erro ao decodificar o payload do webhook"))
	}

	// Ignore the status return because it is not currently used - will be used in messages
	//if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 && len(payload.Entry[0].Changes[0].Value.Statuses) > 0 {
	//	return nil
	//}

	fmt.Print("ID Meta payload : " + payload.Entry[0].Changes[0].Value.MetaData.PhoneNumberId)

	// Check account_meta
	meta, erro := r.metaService.Find(c.Request().Context(), "phone_id="+payload.Entry[0].Changes[0].Value.MetaData.PhoneNumberId)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	// Aqui pegamos o ultimo assistante vinculado a conta meta
	idAssistant, err := r.messageHandler.ValidAssistant(c.Request().Context(), models.WebhookPayload(payload), meta[0])
	if err != nil {
		return responses.Erro(c, http.StatusBadRequest, err)
	}

	// Valida a existência de créditos ou assinatura
	userPlan, err := r.userPlan.Find(c.Request().Context(), "user_id="+meta[0].UserID)
	if err != nil {
		fmt.Println("Erro ao consultar usuario")
		return nil
	}

	fmt.Println("USER PLAN:", userPlan)
	fmt.Println("USER ASASSISTANT:", idAssistant)
	fmt.Println("USER ID:", meta[0].UserID)

	if len(userPlan) == 0 {
		fmt.Println("Usuário sem assinatura ou saldo de créditos")
		return nil
	}

	err = r.messageHandler.Run(c.Request().Context(), models.WebhookPayload(payload), meta[0], idAssistant, meta[0].UserID)
	/*if err != nil {
		return responses.Erro(c, http.StatusBadRequest, err)
	}*/

	return nil

}

// Validate meta parameters
func (r *Webhook) WebhookCheck(c echo.Context) error {
	verifyToken := c.QueryParam("hub.challenge")
	fmt.Println("Start webhook validation process:")

	if verifyToken != "" {
		fmt.Println("Token de verificação encontrado:", verifyToken)
		return c.String(http.StatusOK, verifyToken) //responses.JSON(c, http.StatusOK, verifyToken)
	} else {
		fmt.Println("Token de verificação não encontrado")
		return responses.Erro(c, http.StatusBadRequest, errors.New("Parâmetro hub.verify_token não encontrado na query"))
	}
}

func (r *Webhook) SendMessage(c echo.Context) error {
	messageRequest := new(requests.SendMessageRequest)

	if err := c.Bind(messageRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	// Check account_meta
	meta, erro := r.metaService.Find(c.Request().Context(), "user_id="+messageRequest.PhoneMetaId)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	// Valida a existência de créditos ou assinatura
	userPlan, err := r.userPlan.Find(c.Request().Context(), "user_id="+meta[0].UserID)
	if err != nil {
		fmt.Println("Erro ao consultar usuario")
		return nil
	}

	metaTokens := models.MetaIds{
		PhoneID:    meta[0].PhoneNumberId,
		BusinessId: meta[0].BusinessId,
	}

	fmt.Println("USER PLAN:", userPlan)
	fmt.Println("USER ID:", meta[0].UserID)

	if len(userPlan) == 0 {
		fmt.Println("Usuário sem assinatura ou saldo de créditos")
		return nil
	}

	err = r.messageHandler.SendMessage(c.Request().Context(), messageRequest.CustomerId, messageRequest.Message, messageRequest.ConversationID, metaTokens)
	if err != nil {
		return responses.Erro(c, http.StatusBadRequest, err)
	}

	return nil

}
