package routes

import (
	"go_user_api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	apiV1 := router.Group("/api/v1")
	{
		// router.POST("/signup", handlers.SignUp) is wrong code but correct code:
		apiV1.POST("/signup", handlers.SignUp)
		// router.POST("/Login", handlers.Login) is wrong code but correct code:
		apiV1.POST("/login", handlers.Login)
		// apiV1.GET("/health", handlers.HealthCheck) is wrong code but correct code:
		apiV1.GET("/health", handlers.HealthCheck)
	}

}
