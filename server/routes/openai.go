package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterOpenAiRoutes(e *echo.Echo, openaiController *controllers.OpenAi) {
	e.POST("/assistant", openaiController.Insert)
	// e.GET("/account_meta", metaController.Find)
	// e.GET("/account_meta/:id", metaController.FindId)
	// e.PUT("/account_meta/:id", metaController.Edit)
	// e.DELETE("/account_meta/:id", metaController.Delete)
}
