package services

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"autflow_back/utils"
	"context"
)

type Session struct {
	sessionRepository *repositories.Session
	logger            utils.Logger
}

func NewSession(sessionRepository *repositories.Session, logger utils.Logger) *Session {
	return &Session{
		logger:            logger,
		sessionRepository: sessionRepository,
	}
}

func (r *Session) Insert(ctx context.Context, dt *dto.SessionCreateDTO) (models.Session, error) {
	session := dt.ToSession()

	r.logger.Debugf("Create session: %+v", session)

	idCriado, erro := r.sessionRepository.Insert(ctx, session)
	if erro != nil {
		return models.Session{}, erro
	}

	return idCriado, nil
}

func (r *Session) Find(ctx context.Context, queryString string) ([]models.Session, error) {
	sessions, erro := r.sessionRepository.Find(ctx, queryString)
	if erro != nil {
		return nil, erro
	}

	return sessions, nil
}

func (r *Session) FindId(ctx context.Context, id string) (models.Session, error) {
	session, erro := r.sessionRepository.FindId(ctx, id)
	if erro != nil {
		return models.Session{}, nil

	}

	return *session, nil
}

func (r *Session) Edit(ctx context.Context, id string, dt *dto.SessionCreateDTO) error {
	newSession := dt.ToSession()

	r.logger.Debugf("Edit Session: %+v", newSession)

	return r.sessionRepository.Edit(ctx, id, newSession)
}

func (r *Session) Delete(ctx context.Context, id string) error {
	return r.sessionRepository.Delete(ctx, id)
}

func (r *Session) InsertMessage(ctx context.Context, id string, newMessage models.Message) error {
	return r.sessionRepository.InsertMessage(ctx, id, newMessage)
}

func (r *Session) UpdateSessionField(ctx context.Context, id string, field models.Fields) error {
	return r.sessionRepository.UpdateSessionField(ctx, id, field)
}
