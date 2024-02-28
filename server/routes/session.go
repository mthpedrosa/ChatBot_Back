package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterSessionRoutes(e *echo.Echo, sessionController *controllers.Session) {
	e.POST("/sessions", sessionController.Insert)
	e.GET("/sessions", sessionController.Find)
	e.GET("/sessions/:id", sessionController.FindId)
	e.PUT("/sessions/:id", sessionController.Edit)
	e.PUT("/sessions_fields/:id", sessionController.UpdateSessionField)
	e.DELETE("/sessions/:id", sessionController.Delete)
}
