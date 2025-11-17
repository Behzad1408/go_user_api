package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User defines the structure for a user in the database and API responses.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	// The `json:"-"` tag prevents the password field from ever being exposed in API responses.
	Password string `bson:"password" json:"-"`
}

// Session defines the structure for managing user login state.
// Each document in the sessions collection represents an active login.
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	SessionID string             `bson:"session_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ExpiresAt time.Time          `bson:"expires_at"`
}
