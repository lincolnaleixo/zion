package db

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    Client        *mongo.Client
    LogCollection *mongo.Collection
)

func InitMongoDB(uri string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)
    var err error
    Client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping the database to verify connection
    err = Client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    // Optionally, allow selecting database name from environment variables
    dbName := os.Getenv("MONGODB_DB_NAME")
    if dbName == "" {
        dbName = "smith" // default database name
    }

    LogCollection = Client.Database(dbName).Collection("logs")
    log.Println("Connected to MongoDB!")
}