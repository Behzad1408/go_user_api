package routes

import (
	"go_user_api/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines and groups all the API routes for the application.
func SetupRoutes(router *gin.Engine) {
	// Main group for V1 of the API
	apiV1 := router.Group("/api/v1")
	{
		// Public routes that do not require authentication
		apiV1.POST("/signup", handlers.SignUp)
		apiV1.POST("/login", handlers.Login)
		apiV1.GET("/health", handlers.HealthCheck)

		// Protected routes that require a valid session
		protected := apiV1.Group("/")
		protected.Use(handlers.AuthMiddleware())
		{
			// Route to get the current user's data
			// Full path: GET /api/v1/me
			protected.GET("/me", handlers.GetMyData)

			// Additional protected routes can be added here in the future
			// e.g., protected.PUT("/profile", handlers.UpdateProfile)
		}
	}
}
