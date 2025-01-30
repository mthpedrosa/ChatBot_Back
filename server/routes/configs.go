package routes

import (
	"autflow_back/server/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterConfigRoutes(e *echo.Echo, configController *controllers.ConfigController) {
	e.POST("/config", configController.Create)
	e.GET("/config", configController.GetAll)
	e.GET("/config/:id", configController.GetByID)
	e.PUT("/config/:id", configController.Update)
	e.DELETE("/config/:id", configController.Delete)
}
