package main

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jake-schuler/fll-event-screen/config"
	"github.com/jake-schuler/fll-event-screen/handlers"
)

//go:embed templates/*
var templates embed.FS

func main() {
	// Load environment variables
	if err := godotenv.Load("data/.env"); err != nil {
		panic("Error loading .env file")
	}

	// Initialize database
	db := config.InitDB()

	// Initialize Gin router
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(templates, "templates/*")))
	r.Static("/static", "./static")

	// Setup routes
	handlers.SetupRoutes(r, db)

	// Start server
	if err := r.Run(":8080"); err != nil {
		panic("failed to start server")
	}
}
