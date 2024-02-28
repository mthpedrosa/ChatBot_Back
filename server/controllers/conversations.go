package controllers

import (
	"autflow_back/models/dto"
	"autflow_back/requests"
	"autflow_back/services"
	"autflow_back/src/responses"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Conversation struct {
	conversationService *services.Conversations
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewConversationController(conversation *services.Conversations) *Conversation {
	return &Conversation{
		conversationService: conversation,
	}
}

func (r *Conversation) Insert(c echo.Context) error {
	conversationRequest := new(requests.ConversationRequest)

	if err := c.Bind(conversationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(conversationRequest); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorsMessages := []string{}
		for _, err := range validationErrors {
			errorsMessages = append(errorsMessages, err.Error())
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Validation errors",
			"errors":  errorsMessages,
		})
	}

	dt := &dto.ConversationCreateDTO{
		CustomerId:  conversationRequest.CustomerId,
		Messages:    conversationRequest.Messages,
		WorkflowID:  conversationRequest.WorkflowID,
		OtherFields: conversationRequest.OtherFields,
	}

	createdID, erro := r.conversationService.Insert(c.Request().Context(), dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusCreated, createdID)
}

func (r *Conversation) Find(c echo.Context) error {
	urlParts := strings.Split(c.Request().URL.String(), "?")
	var query string
	if len(urlParts) > 1 {
		query = urlParts[1]
	}

	/*workflow := strings.ToLower(c.Request().URL.Query().Get("workflow"))
	if workflow == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parametro não encontrado"))
	}*/

	conversations, erro := r.conversationService.Find(c.Request().Context(), query)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	// Map goals to a MetaDTO list
	conversationsDTO := make([]dto.ConversationListDTO, len(conversations))
	for i, conversation := range conversations {
		conversationsDTO[i] = dto.ConversationListDTO{
			ID:         conversation.ID,
			CustomerId: conversation.CustomerId,
			WorkflowID: conversation.WorkflowID,
		}
	}

	return responses.JSON(c, http.StatusOK, conversationsDTO)
}

func (r *Conversation) FindId(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parametro não encontrado"))
	}

	conversation, erro := r.conversationService.FindId(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	conversationDTO := dto.ConversationDetailDTO{
		ID:          conversation.ID,
		CustomerId:  conversation.CustomerId,
		Messages:    conversation.Messages,
		WorkflowID:  conversation.WorkflowID,
		CreatedAt:   conversation.CreatedAt,
		UpdateAt:    conversation.UpdateAt,
		OtherFields: conversation.OtherFields,
	}

	return responses.JSON(c, http.StatusOK, conversationDTO)
}

func (r *Conversation) Edit(c echo.Context) error {
	//Extract the ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parametro não encontrado"))
	}

	// Read the new conversation data from the request
	conversationRequest := new(requests.ConversationRequest)

	if err := c.Bind(conversationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(conversationRequest); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorsMessages := []string{}
		for _, err := range validationErrors {
			errorsMessages = append(errorsMessages, err.Error())
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Validation errors",
			"errors":  errorsMessages,
		})
	}

	dt := &dto.ConversationCreateDTO{
		CustomerId:  conversationRequest.CustomerId,
		Messages:    conversationRequest.Messages,
		WorkflowID:  conversationRequest.WorkflowID,
		OtherFields: conversationRequest.OtherFields,
	}

	// Chamar a função para editar a conta
	erro := r.conversationService.Edit(c.Request().Context(), id, dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, "Conta Meta editada com sucesso")
}

func (r *Conversation) Delete(c echo.Context) error {

	//Extract the ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	erro := r.conversationService.Delete(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, "Conversa deletada com sucesso")
}
