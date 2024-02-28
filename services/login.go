package services

import (
	"autflow_back/models"
	"autflow_back/repositories"
	"autflow_back/src/authentication"
	"autflow_back/src/security"
	"autflow_back/utils"
	"context"
	"fmt"
)

type Login struct {
	userRepository *repositories.Users
	logger         utils.Logger
}

func NewLogin(userRepository *repositories.Users, logger utils.Logger) *Login {
	return &Login{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (r *Login) LoginAuth(ctx context.Context, user models.User) (string, error) {
	usuarioSalvoBanco, erro := r.userRepository.FindbyEmail(ctx, user.Email)
	fmt.Println("Usuario salvo no banco ; ", usuarioSalvoBanco)

	if erro = security.CheckPassword(usuarioSalvoBanco.Password, user.Password); erro != nil {
		return "", erro

	}

	token, erro := authentication.CreateToken(usuarioSalvoBanco)
	if erro != nil {
		return "", erro

	}
	return token, nil
}
