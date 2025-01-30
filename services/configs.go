package services

import (
	"autflow_back/models"
	"autflow_back/repositories"
	"context"
)

type ConfigService struct {
	repo *repositories.ConfigRepository
}

func NewConfigService(repo *repositories.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) CreateConfig(ctx context.Context, config models.Config) (string, error) {
	return s.repo.Insert(ctx, config)
}

func (s *ConfigService) UpdateConfig(ctx context.Context, ID string, config models.Config) error {
	return s.repo.Edit(ctx, ID, config)
}

func (s *ConfigService) GetConfigByID(ctx context.Context, ID string) (*models.Config, error) {
	return s.repo.FindByID(ctx, ID)
}

func (s *ConfigService) GetAllConfigs(ctx context.Context, query string) ([]models.Config, error) {
	return s.repo.Find(ctx, query)
}

func (s *ConfigService) DeleteConfig(ctx context.Context, ID string) error {
	return s.repo.Delete(ctx, ID)
}
