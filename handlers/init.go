package handlers

import (
	"backend-template/handlers/authHandler"
	"backend-template/models"

	"github.com/labstack/echo/v4"
)

func All() (routes []models.Route) {

	routes = append(routes, models.Route{
		Path:    "/",
		Method:  echo.GET,
		Handler: health,
	})

	routes = append(routes, authHandler.All("/auth")...)

	return
}
