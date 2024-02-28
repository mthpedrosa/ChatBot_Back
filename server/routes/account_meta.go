package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterMetaRoutes(e *echo.Echo, metaController *controllers.Meta) {
	e.POST("/account_meta", metaController.Insert)
	e.GET("/account_meta", metaController.Find)
	e.GET("/account_meta/:id", metaController.FindId)
	e.PUT("/account_meta/:id", metaController.Edit)
	e.DELETE("/account_meta/:id", metaController.Delete)
}
