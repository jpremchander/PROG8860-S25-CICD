package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"campgrounds-go/config"
	"campgrounds-go/controllers"
	"campgrounds-go/middleware"
	"campgrounds-go/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.ConnectDB()

	// Initialize AWS S3
	config.InitS3()

	// Setup Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("views/*")
	r.Static("/static", "./public")

	// Middleware
	r.Use(middleware.SessionMiddleware())
	r.Use(middleware.FlashMiddleware())

	// Routes
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
