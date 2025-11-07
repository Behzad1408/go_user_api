package user

import (
	"context"
	"log"
	"time"

	"go_user_api/internal/database"
	"go_user_api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FindByEmail(email string) (*models.User, error) {
	var user models.User
	collection := database.GetCollection("mydb", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- بخش دیباگ ---
	log.Printf("Searching for user with email: %s", email)
	// ------------------

	// --- لاگ دیباگ جدید ---
	log.Printf("[FindByEmail] Searching in database: '%s', collection: '%s'", collection.Database().Name(), collection.Name())
	// -----------------------

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		// --- بخش دیباگ ---
		// خطای دقیق را در کنسول چاپ می‌کنیم تا ببینیم چیست
		log.Printf("Error finding user by email: %v", err)
		// ------------------
		return nil, err
	}

	// --- بخش دیباگ ---
	log.Printf("User found: %s", user.Username)
	// ------------------
	return &user, nil
}
