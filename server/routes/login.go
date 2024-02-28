package routes

import (
	"autflow_back/server/controllers"

	"github.com/labstack/echo/v4"
)

/*var routeLogin = Route{
	Uri:                    "/login",
	Method:                 http.MethodPost,
	Function:               controllers.Login,
	RequiresAuthentication: false,
}*/

func RegisterLoginRoutes(e *echo.Echo, loginController *controllers.LoginS) {
	e.POST("/login", loginController.Login)
}
