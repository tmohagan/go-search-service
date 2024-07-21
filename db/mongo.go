package db

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return err
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return err
    }

    Client = client
    log.Println("Connected to MongoDB")
    return nil
}