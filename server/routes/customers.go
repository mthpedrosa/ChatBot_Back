package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterCustomerRoutes(e *echo.Echo, customerController *controllers.Customer) {
	e.POST("/customers", customerController.Insert)
	e.GET("/customers", customerController.Find)
	e.GET("/customers/:id", customerController.FindId)
	e.PUT("/customers/:id", customerController.Edit)
	e.DELETE("/customers/:id", customerController.Delete)
}
