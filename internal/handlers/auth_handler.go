package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" // <-- این import برای خطای duplicate لازمه
	"golang.org/x/crypto/bcrypt"

	"go_user_api/internal/database"
	"go_user_api/internal/models"
	"go_user_api/internal/user" // <-- این import برای FindByEmail لازمه
)

// generateSecureSessionID creates a cryptographically secure, random session ID.
func generateSecureSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func SignUp(c *gin.Context) {
	type SignUpPayload struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var payload SignUpPayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload or missing fields"})
		return
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: payload.Username,
		Email:    payload.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")
	_, err = collection.InsertOne(ctx, &user)

	// --- بخش جدید: مدیریت خطای ایمیل تکراری ---
	if err != nil {
		// این کد چک می‌کنه که آیا خطای برگشتی از MongoDB از نوع "duplicate key" هست یا نه
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "این ایمیل قبلاً ثبت شده است"}) // 409 Conflict
			return
		}
		// اگر خطا از نوع دیگه‌ای بود، خطای عمومی سرور رو برمی‌گردونیم
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد کاربر"})
		return
	}
	// ---------------------------------------------

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// توابع Login, GetMyData, HealthCheck بدون تغییر باقی می‌مانند
func Login(c *gin.Context) {
	type LoginPayload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	foundUser, err := user.FindByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(payload.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	sessionID, err := generateSecureSessionID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	expiresAt := time.Now().Add(24 * 7 * time.Hour)
	session := models.Session{
		ID:        primitive.NewObjectID(),
		SessionID: sessionID,
		UserID:    foundUser.ID,
		ExpiresAt: expiresAt,
	}
	sessionCollection := database.GetCollection("sessions")
	_, err = sessionCollection.InsertOne(ctx, &session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}
	c.SetCookie("session_id", sessionID, 3600*24*7, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": foundUser})
}

func GetMyData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	objID, ok := userID.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var foundUser models.User
	userCollection := database.GetCollection("users")
	err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: User not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": foundUser})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
