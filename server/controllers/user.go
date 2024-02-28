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

type User struct {
	userService *services.User
}

func NewUserController(userService *services.User) *User {
	return &User{
		userService: userService,
	}
}

func (r *User) Insert(c echo.Context) error {
	// Check request body using Bind
	createUserRequest := new(requests.CreateUserRequest)

	if err := c.Bind(createUserRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createUserRequest); err != nil {
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

	dt := &dto.CreateUserDTO{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
		Profile:  createUserRequest.Profile,
	}

	if erro := dt.Prepare("cadastro"); erro != nil {
		return responses.Erro(c, http.StatusBadRequest, erro)

	}

	createdID, erro := r.userService.Insert(c.Request().Context(), dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusCreated, createdID)
}

func (r *User) Find(c echo.Context) error {
	urlParts := strings.Split(c.Request().URL.String(), "?")

	var query string
	if len(urlParts) > 1 {
		query = urlParts[1]
	}

	users, erro := r.userService.Find(c.Request().Context(), query)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	// Map goals to a usersDTO list
	usersDTO := make([]dto.UserListDTO, len(users))
	for i, user := range users {
		usersDTO[i] = dto.UserListDTO{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Profile: user.Profile,
		}
	}

	return responses.JSON(c, http.StatusOK, usersDTO)
}

func (r *User) FindId(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	user, erro := r.userService.FindId(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	userDTO := dto.UserDetailDTO{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Password:   user.Password,
		CreatedAt:  user.CreatedAt,
		UpdateAt:   user.UpdateAt,
		LastActive: user.LastActive,
		Profile:    user.Profile,
	}

	return responses.JSON(c, http.StatusOK, userDTO)
}

func (r *User) Edit(c echo.Context) error {
	//Check if there is permission
	/*if !authentication.HasPermission(c.Request(), config.PermissionsUser) {
		return responses.Erro(c, http.StatusUnauthorized, errors.New("Você não tem permissão para isso"))
	}*/

	// Check request body using Bind
	createUserRequest := new(requests.CreateUserRequest)

	if err := c.Bind(createUserRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createUserRequest); err != nil {
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

	dt := &dto.CreateUserDTO{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
		Profile:  createUserRequest.Profile,
	}

	if erro := dt.Prepare(""); erro != nil {
		return responses.Erro(c, http.StatusBadRequest, erro)
	}

	// Extract the ID from the request
	ID := c.Param("id")
	if ID == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// call the function edit
	erro := r.userService.Edit(c.Request().Context(), ID, dt)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return responses.JSON(c, http.StatusOK, "Usuário editado com sucesso")
}

func (r *User) Delete(c echo.Context) error {
	//Verifica se tem permissao
	/*if !authentication.HasPermission(c.Request(), config.PermissionsUser) {
		return responses.Erro(c, http.StatusUnauthorized, errors.New("Você não tem permissão para isso"))
	}*/

	// Extract the ID from the request
	id := c.Param("id")
	if id == "" {
		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
	}

	// 	// call the function delete
	erro := r.userService.Delete(c.Request().Context(), id)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)

	}

	return responses.JSON(c, http.StatusOK, "Usuário deletado com sucesso")
}
