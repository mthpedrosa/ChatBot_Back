package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterReportsRoutes(e *echo.Echo, reportsController *controllers.Reports) {
	e.POST("/cost_per_user", reportsController.Cost)
}
