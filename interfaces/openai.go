package interfaces

import (
	"autflow_back/models"
	"context"
)

type OpenAIClientRepository interface {
	CreateThread(ctx context.Context) (*models.ThreadResponse, error)
	PostMessage(ctx context.Context, threadID, message string) (string, error)
	StartThreadRun(ctx context.Context, threadID, assistantID string) (string, error)
	GetThreadRunStatus(ctx context.Context, threadID, runID string) (*models.ThreadRun, error)
	GetThreadMessages(ctx context.Context, threadID string) ([]models.MessageThread, error)
	PostToolOutputs(ctx context.Context, threadID, runID, callID string, arrayRespone []models.CallResponse) (string, error)
	ConvertAudioToText(ctx context.Context, filePath string) (string, error)
	CancelRun(ctx context.Context, threadID, runID string) (string, error)
	CreateAssistant(ctx context.Context, dto models.CreateAssistant, model string) (*models.Assistant, error)
	GetAssistant(ctx context.Context, assistantID string) (*models.Assistant, error)
	DeleteAssistant(ctx context.Context, assistantID string) (string, error)
	ListAssistants(ctx context.Context, order string, limit int) ([]models.Assistant, error)
	UpdateAssistant(ctx context.Context, assistantID, model string, dt models.CreateAssistant) (string, error)
}
