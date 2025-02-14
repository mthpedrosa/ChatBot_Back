package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterWebhookRoutes(e *echo.Echo, webhookController *controllers.Webhook) {
	e.POST("/webhook/:id", webhookController.WebhookRun)
	e.GET("/webhook/:id", webhookController.WebhookCheck)
	e.POST("/send-message", webhookController.SendMessage)
}
