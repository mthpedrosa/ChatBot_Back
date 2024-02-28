package services

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"autflow_back/utils"
	"context"
)

type User struct {
	userRepository *repositories.Users
	logger         utils.Logger
}

func NewUser(userRepository *repositories.Users, logger utils.Logger) *User {
	return &User{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (r *User) Insert(ctx context.Context, dt *dto.CreateUserDTO) (string, error) {
	user := dt.ToUser()

	createdID, erro := r.userRepository.Insert(ctx, user)
	if erro != nil {
		return "", erro
	}

	return createdID, nil
}

func (r *User) Find(ctx context.Context, query string) ([]models.User, error) {

	users, erro := r.userRepository.Find(ctx, query)
	if erro != nil {
		return nil, erro

	}

	return users, nil
}

func (r *User) FindId(ctx context.Context, id string) (models.User, error) {

	usuario, erro := r.userRepository.FindId(ctx, id)
	if erro != nil {
		return models.User{}, erro

	}

	return *usuario, nil
}

func (r *User) Edit(ctx context.Context, id string, dt *dto.CreateUserDTO) error {
	newUser := dt.ToUser()
	return r.userRepository.Edit(ctx, id, newUser)
}

func (r *User) Delete(ctx context.Context, id string) error {
	return r.userRepository.Delete(ctx, id)
}
