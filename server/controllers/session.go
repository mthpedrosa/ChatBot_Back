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

type Session struct {
	sessionService *services.Session
}

func NewSessionController(session *services.Session) *Session {
	return &Session{
		sessionService: session,
	}
}

func (r *Session) Insert(c echo.Context) error {
	// Check request body using Bind
	createSessionRequest := new(requests.SessionRequest)

	if err := c.Bind(createSessionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createSessionRequest); err != nil {
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

	dt := &dto.SessionCreateDTO{
		CustomerID:     createSessionRequest.CustomerID,
		WorkflowId:     createSessionRequest.WorkflowId,
		ConversationId: createSessionRequest.ConversationId,
		Status:         createSessionRequest.Status,
		Tags:           createSessionRequest.Tags,
		OtherFields:    createSessionRequest.OtherFields,
		Messages:       createSessionRequest.Messages,
		LastNode:       createSessionRequest.LastNode,
	}

	if erro := dt.Validate(); erro != nil {
		return responses.Erro(c, http.StatusBadRequest, erro)

	}

	session, erro := r.sessionService.Insert(c.Request().Context(), dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusCreated, session)
}

func (r *Session) Find(c echo.Context) error {

	// Extrair a query string diretamente da requisição
	queryString := c.Request().URL.RawQuery

	sessions, erro := r.sessionService.Find(c.Request().Context(), queryString)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	// Map goals to a customerDTO list
	sessionDTO := make([]dto.SessionListDTO, len(sessions))
	for i, session := range sessions {
		sessionDTO[i] = dto.SessionListDTO{
			ID:             session.ID,
			CustomerID:     session.CustomerID,
			WorkflowId:     session.WorkflowId,
			ConversationId: session.ConversationId,
			Status:         session.Status,
		}
	}

	return responses.JSON(c, http.StatusOK, sessionDTO)
}

func (r *Session) FindId(c echo.Context) error {
	// Extract the session ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	session, erro := r.sessionService.FindId(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	sessionDTO := dto.SessionDetailDTO{
		ID:             session.ID,
		CustomerID:     session.CustomerID,
		WorkflowId:     session.WorkflowId,
		ConversationId: session.ConversationId,
		Status:         session.Status,
		Duration:       session.Duration,
		CreatedAt:      session.CreatedAt,
		UpdateAt:       session.UpdateAt,
		FinishedAt:     session.FinishedAt,
		Tags:           session.Tags,
		OtherFields:    session.OtherFields,
		Messages:       session.Messages,
		LastNode:       session.Status,
	}

	return responses.JSON(c, http.StatusOK, sessionDTO)
}

func (r *Session) Edit(c echo.Context) error {
	// Check request body using Bind
	createSessionRequest := new(requests.SessionRequest)

	if err := c.Bind(createSessionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createSessionRequest); err != nil {
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

	dt := &dto.SessionCreateDTO{
		CustomerID:     createSessionRequest.CustomerID,
		WorkflowId:     createSessionRequest.WorkflowId,
		ConversationId: createSessionRequest.ConversationId,
		Status:         createSessionRequest.Status,
		Tags:           createSessionRequest.Tags,
		OtherFields:    createSessionRequest.OtherFields,
		Messages:       createSessionRequest.Messages,
		LastNode:       createSessionRequest.LastNode,
	}
	if erro := dt.Validate(); erro != nil {
		return responses.Erro(c, http.StatusBadRequest, erro)

	}

	// Extract the customer ID from the request
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// Call the function to edit the customer
	erro := r.sessionService.Edit(c.Request().Context(), ID, dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusOK, "Session editada com sucesso")
}

func (r *Session) Delete(c echo.Context) error {

	// Extract the session ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// Call the function to delete the session
	erro := r.sessionService.Delete(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, "Session deletado com sucesso")
}

func (r *Session) UpdateSessionField(c echo.Context) error {
	// Check request body using Bind
	otherFieldRequest := new(requests.SessionOtherRequest)

	if err := c.Bind(otherFieldRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(otherFieldRequest); err != nil {
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

	// Extract the customer ID from the request
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// Call the function to edit the customer
	erro := r.sessionService.UpdateSessionField(c.Request().Context(), ID, otherFieldRequest.OtherFields)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusOK, "Session editada com sucesso")
}

// parseQueryString analisa a query string e retorna um mapa dos parâmetros
func parseQueryString(query string) map[string]string {
	params := make(map[string]string)
	for _, param := range strings.Split(query, "&") {
		pair := strings.SplitN(param, "=", 2)
		if len(pair) == 2 {
			params[pair[0]] = pair[1]
		}
	}
	return params
}

// Função para analisar a string de other_fields e retornar um mapa de filtros
func parseOtherFields(otherFields string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(otherFields, ";")
	for _, pair := range pairs {
		keyValue := strings.Split(pair, "=")
		if len(keyValue) == 2 {
			result[keyValue[0]] = keyValue[1]
		}
	}
	return result
}
