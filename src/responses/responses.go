package responses

import "github.com/labstack/echo/v4"

// JSON retorna uma resposta JSON para a requisição no Echo.
func JSON(c echo.Context, statusCode int, dados interface{}) error {
	return c.JSON(statusCode, dados)
}

// Erro retorna uma resposta JSON de erro para a requisição no Echo.
func Erro(c echo.Context, statusCode int, erro error) error {
	return JSON(c, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
}
