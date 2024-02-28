package controllers

import (
	"autflow_back/models"
	"autflow_back/services"
	"autflow_back/src/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginS struct {
	loginService *services.Login
}

func NewLoginController(loginService *services.Login) *LoginS {
	return &LoginS{
		loginService: loginService,
	}
}

func (r *LoginS) Login(c echo.Context) error {

	corpoRequisicao, erro := ioutil.ReadAll(c.Request().Body)
	if erro != nil {
		return responses.Erro(c, http.StatusUnprocessableEntity, erro)

	}

	var user models.User
	if erro = json.Unmarshal(corpoRequisicao, &user); erro != nil {
		return responses.Erro(c, http.StatusBadRequest, erro)

	}

	token, erro := r.loginService.LoginAuth(c.Request().Context(), user)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	return c.String(http.StatusOK, token)
}
