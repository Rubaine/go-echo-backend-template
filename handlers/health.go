package handlers

import (
	config "backend-template/config"
	models "backend-template/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func health(c echo.Context) error {

	return c.JSON(http.StatusOK, models.HealthModel{
		Status:  "Healthy !",
		Version: config.Version,
	})
}
