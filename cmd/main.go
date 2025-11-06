package main

import (
	"log"

	"go_user_api/internal/database"
	"go_user_api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 1. اتصال به MongoDB
	if err := database.Connect(); err != nil {
		log.Fatal("خطا در اتصال به دیتابیس:", err)
	}

	// 2. ساخت router از Gin
	router := gin.Default()

	// 3. تنظیم routeها
	routes.SetupRoutes(router)

	// 4. اجرای سرور روی پورت 8080
	log.Println("سرور در حال اجرا روی پورت 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("خطا در اجرای سرور:", err)
	}
}
