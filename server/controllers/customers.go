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

type Customer struct {
	customerService *services.Customer
}

func NewCustomerController(customer *services.Customer) *Customer {
	return &Customer{
		customerService: customer,
	}
}

func (r *Customer) Insert(c echo.Context) error {
	// Check request body using Bind
	createMetaRequest := new(requests.CustomerRequest)

	if err := c.Bind(createMetaRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createMetaRequest); err != nil {
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

	dt := &dto.CreateCustomerDTO{
		Name:        createMetaRequest.Name,
		Email:       createMetaRequest.Email,
		Phone:       createMetaRequest.Phone,
		WhatsAppID:  createMetaRequest.WhatsAppID,
		OtherFields: createMetaRequest.OtherFields,
	}

	if erro := dt.Validate(); erro != nil {
		return responses.Erro(c, http.StatusBadRequest, erro)

	}

	idCriado, erro := r.customerService.Insert(c.Request().Context(), dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	response := map[string]string{"id": idCriado}
	return responses.JSON(c, http.StatusCreated, response)
}

func (r *Customer) Find(c echo.Context) error {

	urlParts := strings.Split(c.Request().URL.String(), "?")
	var query string
	if len(urlParts) > 1 {
		query = urlParts[1]
	}

	customers, erro := r.customerService.Find(c.Request().Context(), query)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	// Map goals to a customerDTO list
	customerDTO := make([]dto.CustomerListDTO, len(customers))
	for i, customer := range customers {
		customerDTO[i] = dto.CustomerListDTO{
			ID:          customer.ID,
			Name:        customer.Name,
			Email:       customer.Email,
			WhatsAppID:  customer.WhatsAppID,
			OtherFields: customer.OtherFields,
		}
	}

	return responses.JSON(c, http.StatusOK, customerDTO)
}

func (r *Customer) FindId(c echo.Context) error {
	// Extract the customer ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	customer, erro := r.customerService.FindId(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	customerDTO := dto.CustomerDetailDTO{
		ID:          customer.ID,
		Name:        customer.Name,
		Email:       customer.Email,
		Phone:       customer.Phone,
		WhatsAppID:  customer.WhatsAppID,
		CreatedAt:   customer.CreatedAt,
		UpdateAt:    customer.UpdateAt,
		OtherFields: customer.OtherFields,
	}

	return responses.JSON(c, http.StatusOK, customerDTO)
}

func (r *Customer) Edit(c echo.Context) error {
	// Check request body using Bind
	createMetaRequest := new(requests.CustomerRequest)

	if err := c.Bind(createMetaRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createMetaRequest); err != nil {
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

	dt := &dto.CreateCustomerDTO{
		Name:        createMetaRequest.Name,
		Email:       createMetaRequest.Email,
		Phone:       createMetaRequest.Phone,
		WhatsAppID:  createMetaRequest.WhatsAppID,
		OtherFields: createMetaRequest.OtherFields,
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
	erro := r.customerService.Edit(c.Request().Context(), ID, dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusOK, "Cliente editado com sucesso")
}

func (r *Customer) Delete(c echo.Context) error {

	// Extract the customer ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// Call the function to delete the customer
	erro := r.customerService.Delete(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, "Cliente deletado com sucesso")
}
