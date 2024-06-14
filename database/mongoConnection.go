package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client = DBinstance()

func DBinstance() *mongo.Client {
	fmt.Println("MongoDB connection in progress...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
	MongoDB := os.Getenv("MONGODB_URL")
	// Set up a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			fmt.Printf("Started command: %v\n", evt.Command)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			fmt.Printf("Succeeded command: %v\n", evt.Reply)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			fmt.Printf("Failed command: %v\n", evt.Failure)
		},
	}
	clientOptions := options.Client().ApplyURI(MongoDB).SetMonitor(monitor)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Check the connection
    err = client.Ping(context.TODO(), readpref.Primary())
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
	fmt.Println("MongoDB is connected successfully")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	return collection
}
