package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() error { // ← error اضافه شد
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//connectionString := "mongodb://myuser:mypass@localhost:27017"
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err // ← error برگردون
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err // ← error برگردون
	}

	Client = client
	log.Println("Connected to MongoDB!")

	return nil // ← اگه موفق بود، nil برگردون
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return Client.Database(databaseName).Collection(collectionName)
}
