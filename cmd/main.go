package main

import (
	"log"

	"go_user_api/internal/database"
	"go_user_api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// 1. Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect to the MongoDB database
	if err := database.Connect(); err != nil {
		log.Fatal("خطا در اتصال به دیتابیس:", err)
	}

	// 3. Create a new Gin router with default middleware
	router := gin.Default()

	// 4. Set up all API routes
	routes.SetupRoutes(router)

	// 5. Start the HTTP server and listen on port 8080
	log.Println("سرور در حال اجرا روی پورت 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("خطا در اجرای سرور:", err)
	}
}
