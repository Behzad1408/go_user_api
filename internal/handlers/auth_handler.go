package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go_user_api/internal/database"
	"go_user_api/internal/models"
)

func SignUp(c *gin.Context) {
	// 1. یک متغیر User بساز
	var user models.User

	// 2. داده‌های JSON رو از request بگیر
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های نامعتبر"})
		return
	}

	// 3. چک کن فیلدها پر هستن
	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "همه فیلدها الزامی هستند"})
		return
	}

	// 4. Password رو hash کن
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در هش کردن رمز عبور"})
		return
	}
	user.Password = string(hashedPassword)

	// 5. ID جدید بساز
	user.ID = primitive.NewObjectID()

	// 6. context بساز
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 7. collection رو بگیر
	collection := database.GetCollection("mydb", "users")

	// 8. User رو در MongoDB ذخیره کن
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ذخیره‌سازی کاربر"})
		return
	}

	// 9. جواب موفقیت برگردون
	c.JSON(http.StatusCreated, gin.H{
		"message": "ثبت‌نام با موفقیت انجام شد",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
