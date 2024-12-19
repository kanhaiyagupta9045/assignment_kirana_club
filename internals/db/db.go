package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	once        sync.Once
)

func DBConnection() *mongo.Client {
	once.Do(func() {
		var err error
		mongodb_url := os.Getenv("MONGODB_URL")
		if mongodb_url == "" {
			log.Fatal("please provide the mongodb connection string")
		}
		MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongodb_url))
		if err != nil {
			log.Fatalf("%v", err)
		}
		if err := MongoClient.Ping(context.TODO(), nil); err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println("Database connected successfully")
	})

	return MongoClient
}
