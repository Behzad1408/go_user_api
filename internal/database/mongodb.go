package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DbName string // Will hold the database name from .env

// createUniqueIndexes ensures that the necessary unique indexes are created on collections.
func createUniqueIndexes(ctx context.Context) {
	userCollection := Client.Database(DbName).Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	indexName, err := userCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("Error creating unique index for email: %v", err)
	}

	log.Printf("Unique index '%s' on collection 'users' in database '%s' ensured.", indexName, DbName)
}

// Connect establishes a connection to MongoDB and creates necessary indexes.
func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	DbName = os.Getenv("MONGO_DB_NAME")
	if DbName == "" {
		log.Fatal("MONGO_DB_NAME environment variable not set")
	}

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	log.Printf("Successfully connected to MongoDB! Using database: %s", DbName)

	indexCtx, indexCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer indexCancel()
	createUniqueIndexes(indexCtx)

	return nil
}

// GetCollection is a helper function to retrieve a collection from the database.
// It now uses the global DbName variable and only needs the collection name.
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database(DbName).Collection(collectionName)
}
