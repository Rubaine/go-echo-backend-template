package main

import (
	"embed"
	"fmt"
	"os"

	"example.com/template/config"
	_ "example.com/template/docs"
	"example.com/template/handlers"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed public/*
var Folder embed.FS

// @title Echo Backend Template API
// @version 1.0
// @description This is a template for a Go Echo backend API.
// @termsOfService https://example.com/terms

// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email
func main() {

	// Initialize echo
	api := echo.New()
	api.HideBanner = true
	api.HTTPErrorHandler = handlers.OnError

	origins := []string{"https://site.fr"}
	if log.GetLevel() == log.DebugLevel {
		origins = append(origins, "*") // Allow all origins in debug mode
	}

	crs := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers", handlers.TokenKeyName},
		AllowMethods:     []string{echo.POST, echo.GET, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	})
	api.Use(crs)

	api.Use(
		middleware.BodyLimit(config.Config.BodySizeLimit),
		middleware.Gzip(),
		handlers.AuthMiddleware,
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `[${time_rfc3339}] ${status} ${method} ${uri} ${latency_human} ${bytes_in} ${bytes_out} ${remote_ip}` + "\n",
			Output: os.Stdout,
		}),
	)

	// Register api routes
	for _, handler := range handlers.All() {
		api.Add(handler.Method, handler.Path, handler.Handler, handler.Middlewares...)
	}

	// Swagger
	if log.GetLevel() == log.DebugLevel {
		api.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Start public API
	err := api.Start(fmt.Sprintf(":%s", config.Config.ListenPort))
	if err != nil {
		log.Fatal("Public API handler stopped", "error", err)
	}
}
