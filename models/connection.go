package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment variables")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	MongoDB := os.Getenv("MONGODB")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database...")
	return client
}

var Client *mongo.Client = connectDB()

func UserCollection(client *mongo.Client) *mongo.Collection {
	var collection *mongo.Collection = client.Database("goauthDB").Collection("User")
	return collection
}
func TokenCollection(client *mongo.Client) *mongo.Collection {
	var collection *mongo.Collection = client.Database("goauthDB").Collection("Token")
	return collection
}
