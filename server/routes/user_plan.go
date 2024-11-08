package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterUserPlanRoutes(e *echo.Echo, userPlanController *controllers.UserPlanController) {
	e.POST("/user_plan", userPlanController.Insert)
	e.GET("/user_plan/:id", userPlanController.FindId)
	e.PUT("/user_plan/:id", userPlanController.Edit)
	e.DELETE("/user_plan/:id", userPlanController.Delete)
}
