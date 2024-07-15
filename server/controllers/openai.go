package controllers

import (
	"autflow_back/models/dto"
	"autflow_back/requests"
	"autflow_back/services"
	"autflow_back/src/responses"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type OpenAi struct {
	openaiService *services.OpenAi
}

func NewOpenAiController(openaiService *services.OpenAi) *OpenAi {
	return &OpenAi{
		openaiService: openaiService,
	}
}

func (o *OpenAi) Insert(c echo.Context) error {
	// Check request body using Bind
	createAssistantRequest := new(requests.CreateAssistantRequest)

	if err := c.Bind(createAssistantRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createAssistantRequest); err != nil {
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

	// ID of the user creating the user
	// creatorUser, erro := authentication.ExtractIdToken(c.Request())
	// if erro != nil {
	// 	return responses.Erro(c, http.StatusBadRequest, erro)
	// }

	dt := &dto.CreateAssistantDTO{
		Name:         createAssistantRequest.Name,
		Instructions: createAssistantRequest.Instructions,
	}

	// Call the service with the Meta
	createdID, erro := o.openaiService.Insert(c.Request().Context(), dt, createAssistantRequest.IdCustomer)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusCreated, createdID)
}
func (o *OpenAi) Edit(c echo.Context) error {
	// Check request body using Bind
	createAssistantRequest := new(requests.CreateAssistantRequest)

	if err := c.Bind(createAssistantRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createAssistantRequest); err != nil {
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

	// Extract the account ID from the request (e.g., from a parameter in the URL)
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	dt := &dto.CreateAssistantDTO{
		Name:         createAssistantRequest.Name,
		Instructions: createAssistantRequest.Instructions,
	}

	retorno, erro := o.openaiService.Edit(c.Request().Context(), dt, ID)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusOK, "Assistant editado com sucesso: "+retorno)
}

func (o *OpenAi) FindAll(c echo.Context) error {
	order := c.QueryParam("order")
	limit := c.QueryParam("limit")
	num := 20

	if order == "" {
		order = "desc"
	}
	if limit != "" {
		numConvert, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		num = numConvert
	}

	openai, erro := o.openaiService.FindAll(c.Request().Context(), order, num)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, openai)
}

func (o *OpenAi) FindId(c echo.Context) error {
	// Extract the user ID from the request (e.g., from a parameter in the URL)
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	openai, erro := o.openaiService.FindId(c.Request().Context(), ID)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, openai)
}

func (o *OpenAi) Delete(c echo.Context) error {
	// Extract the account ID from the request (e.g., from a parameter in the URL)
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	// delete the account
	retorno, erro := o.openaiService.Delete(c.Request().Context(), ID)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}
	return responses.JSON(c, http.StatusOK, retorno)
}
