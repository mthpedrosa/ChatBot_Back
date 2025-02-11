package interfaces

import (
	"autflow_back/models"
	"context"
)

type WhatsappRepository interface {
	InteractiveMessage(ctx context.Context, text string, buttonsArray []models.Button, customer models.Customer, meta models.MetaIds) error
	SimpleMessage(ctx context.Context, messageSend string, customer models.Customer, meta models.MetaIds) error
	InteractiveListMessage(ctx context.Context, customer models.Customer, meta models.MetaIds) error
	GetUrlMedia(ctx context.Context, mediaID string) (string, error)
	DownloadMedia(ctx context.Context, url, name string) (string, error)
	InteractiveMessageList(ctx context.Context, customer models.Customer, meta models.MetaIds, bodyText string, rows []models.Row) error
	ContactMessage(ctx context.Context, customer models.Customer, meta models.MetaIds, name, phone string) error
}
