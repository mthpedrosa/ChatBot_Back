package controllers

import (
	"autflow_back/models"
	"autflow_back/services"
	"context"
	"github.com/labstack/echo/v4"

	"net/http"
)

type ConfigController struct {
	service *services.ConfigService
}

func NewConfigController(service *services.ConfigService) *ConfigController {
	return &ConfigController{service: service}
}

func (cc *ConfigController) Create(c echo.Context) error {
	var config models.Config
	if err := c.Bind(&config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	id, err := cc.service.CreateConfig(context.Background(), config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func (cc *ConfigController) GetAll(c echo.Context) error {
	configs, err := cc.service.GetAllConfigs(context.Background(), "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, configs)
}

func (cc *ConfigController) GetByID(c echo.Context) error {
	id := c.Param("id")
	config, err := cc.service.GetConfigByID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, config)
}

func (cc *ConfigController) Update(c echo.Context) error {
	id := c.Param("id")
	var config models.Config

	if err := c.Bind(&config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := cc.service.UpdateConfig(context.Background(), id, config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Config updated successfully"})
}

func (cc *ConfigController) Delete(c echo.Context) error {
	id := c.Param("id")
	err := cc.service.DeleteConfig(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Config deleted successfully"})
}
