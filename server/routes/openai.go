package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterOpenAiRoutes(e *echo.Echo, openaiController *controllers.OpenAi) {
	e.POST("/assistant", openaiController.Insert)
	e.GET("/assistant", openaiController.FindAll)
	e.GET("/assistant/:id", openaiController.FindId)
	e.PUT("/assistant/:id", openaiController.Edit)
	e.DELETE("/assistant/:id", openaiController.Delete)
	e.GET("/assistant_user/:id", openaiController.FindAllUser)
}
