package db

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/bson"
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

func PerformSearch(query string) ([]bson.M, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    database := Client.Database("your_database_name") // Make sure this is correct
    postsCollection := database.Collection("posts")
    projectsCollection := database.Collection("projects")

    log.Printf("Searching for query: %s", query)

    filter := bson.M{
        "$or": []bson.M{
            {"title": bson.M{"$regex": query, "$options": "i"}},
            {"summary": bson.M{"$regex": query, "$options": "i"}},
            {"content": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    var results []bson.M

    // Search Posts
    postsCursor, err := postsCollection.Find(ctx, filter)
    if err != nil {
        log.Printf("Error searching posts: %v", err)
        return nil, err
    }
    defer postsCursor.Close(ctx)

    for postsCursor.Next(ctx) {
        var result bson.M
        err := postsCursor.Decode(&result)
        if err != nil {
            log.Printf("Error decoding post result: %v", err)
            return nil, err
        }
        result["type"] = "post"
        results = append(results, result)
    }

    log.Printf("Found %d posts", len(results))

    // Search Projects
    projectsCursor, err := projectsCollection.Find(ctx, filter)
    if err != nil {
        log.Printf("Error searching projects: %v", err)
        return nil, err
    }
    defer projectsCursor.Close(ctx)

    for projectsCursor.Next(ctx) {
        var result bson.M
        err := projectsCursor.Decode(&result)
        if err != nil {
            log.Printf("Error decoding project result: %v", err)
            return nil, err
        }
        result["type"] = "project"
        results = append(results, result)
    }

    log.Printf("Found %d total results", len(results))

    return results, nil
}