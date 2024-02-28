package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterUsersRoutes(e *echo.Echo, userController *controllers.User) {
	e.POST("/users", userController.Insert)
	e.GET("/users", userController.Find)
	e.GET("/users/:id", userController.FindId)
	e.PUT("/users/:id", userController.Edit)
	e.DELETE("/users/:id", userController.Delete)
}
