package services

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"autflow_back/utils"
	"context"
)

type Meta struct {
	metaRepository *repositories.Metas
	logger         utils.Logger
}

func NewMeta(
	metaRepository *repositories.Metas,
	logger utils.Logger) *Meta {
	return &Meta{

		metaRepository: metaRepository,
		logger:         logger,
	}
}

func (r *Meta) Insert(ctx context.Context, dt *dto.CreateMetaDTO) (string, error) {
	meta := dt.ToMeta()

	r.logger.Debugf("Account Meta: %+v", meta)

	idCriado, erro := r.metaRepository.Insert(ctx, meta)
	if erro != nil {
		return "", erro
	}

	return idCriado, nil
}

func (r *Meta) Find(ctx context.Context, query string) ([]models.Meta, error) {
	metas, erro := r.metaRepository.Find(ctx, query)
	if erro != nil {
		return nil, erro
	}

	return metas, nil
}

func (r *Meta) FindId(ctx context.Context, id string) (models.Meta, error) {

	meta, erro := r.metaRepository.FindId(ctx, id)
	if erro != nil {
		return models.Meta{}, erro
	}

	return *meta, nil
}

func (r *Meta) Edit(ctx context.Context, dt *dto.CreateMetaDTO, id string) error {
	meta := dt.ToMeta()

	return r.metaRepository.Edit(ctx, id, meta)
}

func (r *Meta) Delete(ctx context.Context, id string) error {
	return r.metaRepository.Delete(ctx, id)
}
