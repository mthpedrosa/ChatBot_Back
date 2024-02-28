package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterConversationsRoutes(e *echo.Echo, conversationsController *controllers.Conversation) {
	e.POST("/conversation", conversationsController.Insert)
	e.GET("/conversation", conversationsController.Find)
	e.GET("/conversation/:id", conversationsController.FindId)
	e.PUT("/conversation/:id", conversationsController.Edit)
	e.DELETE("/conversation/:id", conversationsController.Delete)
}
