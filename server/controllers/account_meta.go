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

type Meta struct {
	metaService *services.Meta
}

//var validate = validator.New(validator.WithRequiredStructEnabled())

func NewMetaController(metaService *services.Meta) *Meta {
	return &Meta{
		metaService: metaService,
	}
}

func (r *Meta) Insert(c echo.Context) error {
	// Check request body using Bind
	createMetaRequest := new(requests.CreateMetaRequest)

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

	// ID of the user creating the user
	// creatorUser, erro := authentication.ExtractIdToken(c.Request())
	// if erro != nil {
	// 	return responses.Erro(c, http.StatusBadRequest, erro)
	// }

	dt := &dto.CreateMetaDTO{
		Name:          createMetaRequest.Name,
		PhoneNumberId: createMetaRequest.PhoneNumberId,
		BusinessId:    createMetaRequest.BusinessId,
		Assistants:    createMetaRequest.Assistants,
		UserID:        createMetaRequest.UserID,
	}

	// Call the service with the Meta
	createdID, erro := r.metaService.Insert(c.Request().Context(), dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusCreated, createdID)
}

func (r *Meta) Find(c echo.Context) error {
	urlParts := strings.Split(c.Request().URL.String(), "?")
	var query string
	if len(urlParts) > 1 {
		query = urlParts[1]
	}

	metas, erro := r.metaService.Find(c.Request().Context(), query)
	if erro != nil {

		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	// Map goals to a MetaDTO list
	metasDTO := make([]dto.MetaListDTO, len(metas))
	for i, meta := range metas {
		metasDTO[i] = dto.MetaListDTO{
			ID:            meta.ID,
			Name:          meta.Name,
			PhoneNumberId: meta.PhoneNumberId,
			BusinessId:    meta.BusinessId,
			UserID:        meta.UserID,
		}
	}

	return responses.JSON(c, http.StatusOK, metasDTO)
}

func (r *Meta) FindId(c echo.Context) error {
	// Extract the user ID from the request (e.g., from a parameter in the URL)
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	meta, erro := r.metaService.FindId(c.Request().Context(), ID)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	metaDTO := dto.MetaDetailDTO{
		ID:            meta.ID,
		Name:          meta.Name,
		PhoneNumberId: meta.PhoneNumberId,
		BusinessId:    meta.BusinessId,
		CreatedAt:     meta.CreatedAt,
		UpdateAt:      meta.UpdateAt,
		Assistants:    meta.Assistants,
		UserID:        meta.UserID,
	}

	return responses.JSON(c, http.StatusOK, metaDTO)
}

func (r *Meta) Edit(c echo.Context) error {
	// Check request body using Bind
	createMetaRequest := new(requests.CreateMetaRequest)

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

	// Extract the account ID from the request (e.g., from a parameter in the URL)
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	dt := &dto.CreateMetaDTO{
		Name:          createMetaRequest.Name,
		PhoneNumberId: createMetaRequest.PhoneNumberId,
		BusinessId:    createMetaRequest.BusinessId,
		Assistants:    createMetaRequest.Assistants,
		UserID:        createMetaRequest.UserID,
	}

	erro := r.metaService.Edit(c.Request().Context(), dt, ID)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusOK, "Conta Meta editada com sucesso")
}

func (r *Meta) Delete(c echo.Context) error {
	// Extract the account ID from the request (e.g., from a parameter in the URL)
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	// delete the account
	erro := r.metaService.Delete(c.Request().Context(), ID)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}
	return responses.JSON(c, http.StatusOK, "Conta Meta deletada com sucesso")
}
