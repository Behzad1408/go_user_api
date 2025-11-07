package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go_user_api/internal/database"
	"go_user_api/internal/models"
	"go_user_api/internal/user"
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

	// --- لاگ دیباگ جدید ---
	log.Printf("[SignUp] Inserting user into database: '%s', collection: '%s'", collection.Database().Name(), collection.Name())
	// -----------------------

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

// در انتهای فایل auth_handler.go

func Login(c *gin.Context) {
	// یک ساختار موقت برای گرفتن ورودی‌های لاگین
	type LoginPayload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var payload LoginPayload

	// بایند کردن JSON ورودی
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// مرحله ۲.۱: پیدا کردن کاربر در دیتابیس
	// ما اینجا از پکیج user که هنوز نساختیم استفاده می‌کنیم
	foundUser, err := user.FindByEmail(payload.Email)
	if err != nil {
		// این خطا یعنی یا کاربر پیدا نشد یا خطای دیتابیس رخ داده
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid email or password"})
		return
	}

	// مرحله ۲.۲: مقایسه رمز عبور ورودی با رمز عبور هش شده در دیتابیس
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(payload.Password))
	if err != nil {
		// اگر err != nil باشد، یعنی رمزها مطابقت ندارند
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid email or password"})
		return
	}

	// مرحله ۳: پاسخ موفقیت‌آمیز
	// اگر کد به اینجا برسد، یعنی کاربر پیدا شده و رمز عبور هم درست بوده است
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successful"})
}
