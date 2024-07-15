package interfaces

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"context"
)

type OpenAIClientRepository interface {
	CreateThread(ctx context.Context) (*models.ThreadResponse, error)
	PostMessage(ctx context.Context, threadID, message string) (string, error)
	StartThreadRun(ctx context.Context, threadID string) (string, error)
	GetThreadRunStatus(ctx context.Context, threadID, runID string) (*models.ThreadRun, error)
	GetThreadMessages(ctx context.Context, threadID string) ([]models.MessageThread, error)
	PostToolOutputs(ctx context.Context, threadID, runID, callID string, arrayRespone []models.CallResponse) (string, error)
	ConvertAudioToText(ctx context.Context, filePath string) (string, error)
	CancelRun(ctx context.Context, threadID, runID string) (string, error)
	CreateAssistant(ctx context.Context, dto dto.CreateAssistantDTO, model string) (*models.Assistant, error)
	GetAssistant(ctx context.Context, assistantID string) (*models.Assistant, error)
	DeleteAssistant(ctx context.Context, assistantID string) (string, error)
}
