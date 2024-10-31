package controllers

import (
	"autflow_back/requests"
	"autflow_back/services"
	"autflow_back/src/responses"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Reports struct {
	reportsService *services.Reports
}

func NewReportsController(reportsService *services.Reports) *Reports {
	return &Reports{
		reportsService: reportsService,
	}
}

func (o *Reports) Cost(c echo.Context) error {
	costRequest := new(requests.CostParams)

	if err := c.Bind(costRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(costRequest); err != nil {
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

	fmt.Print(costRequest)

	createdID, erro := o.reportsService.Cost(c.Request().Context(), costRequest)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusCreated, createdID)
}
