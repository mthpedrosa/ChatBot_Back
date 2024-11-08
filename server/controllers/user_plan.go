package controllers

import (
	"autflow_back/requests"
	"autflow_back/services"
	"autflow_back/src/responses"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserPlanController struct {
	userPlanService *services.UserPlanService
}

// NewUserPlanController cria uma nova instância de UserPlanController
func NewUserPlanController(userPlanService *services.UserPlanService) *UserPlanController {
	return &UserPlanController{userPlanService: userPlanService}
}

// CreateUserPlan cria um novo plano de usuário
func (ctrl *UserPlanController) Insert(c echo.Context) error {
	userPlan := new(requests.UserPlanRequest)

	if err := c.Bind(userPlan); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(userPlan); err != nil {
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

	userPlanId, err := ctrl.userPlanService.Insert(c.Request().Context(), *userPlan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create user plan"})
	}

	return c.JSON(http.StatusCreated, userPlanId)
}

// EditUserPlan atualiza um plano de usuário existente
func (ctrl *UserPlanController) Edit(c echo.Context) error {
	// Check request body using Bind
	userPlanRequest := new(requests.UserPlanRequest)

	if err := c.Bind(userPlanRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(userPlanRequest); err != nil {
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

	if err := ctrl.userPlanService.Edit(c.Request().Context(), ID, *userPlanRequest); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update user plan"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "user plan updated successfully"})
}

// DeleteUserPlan exclui um plano de usuário existente
func (ctrl *UserPlanController) Delete(c echo.Context) error {
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	if err := ctrl.userPlanService.Delete(c.Request().Context(), ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete user plan"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "user plan deleted successfully"})
}

// GetUserPlan retorna um plano de usuário pelo ID
func (ctrl *UserPlanController) FindId(c echo.Context) error {
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("parametro não encontrado"))
	}

	userPlan, err := ctrl.userPlanService.FindId(c.Request().Context(), ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user plan not found"})
	}

	return c.JSON(http.StatusOK, userPlan)
}
