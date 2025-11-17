package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"go_user_api/internal/database"
	"go_user_api/internal/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var session models.Session
		// FIX: No longer passing "mydb"
		sessionCollection := database.GetCollection("sessions")

		err = sessionCollection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			return
		}

		if time.Now().After(session.ExpiresAt) {
			_, _ = sessionCollection.DeleteOne(ctx, bson.M{"_id": session.ID})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again."})
			return
		}

		c.Set("userID", session.UserID)
		c.Next()
	}
}
