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

func PerformSearch(query string, page, resultsPerPage int64) ([]bson.M, int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    database := Client.Database("test")
    postsCollection := database.Collection("posts")
    projectsCollection := database.Collection("projects")

    log.Printf("Searching for query: %s, page: %d, resultsPerPage: %d", query, page, resultsPerPage)

    filter := bson.M{
        "$or": []bson.M{
            {"title": bson.M{"$regex": query, "$options": "i"}},
            {"summary": bson.M{"$regex": query, "$options": "i"}},
            {"content": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    skip := (page - 1) * resultsPerPage
    findOptions := options.Find().
        SetSkip(skip).
        SetLimit(resultsPerPage)

    var results []bson.M
    var totalCount int64

    // Search Posts
    postsCursor, err := postsCollection.Find(ctx, filter, findOptions)
    if err != nil {
        log.Printf("Error searching posts: %v", err)
        return nil, 0, err
    }
    defer postsCursor.Close(ctx)

    for postsCursor.Next(ctx) {
        var result bson.M
        if err := postsCursor.Decode(&result); err != nil {
            log.Printf("Error decoding post result: %v", err)
            return nil, 0, err
        }
        result["type"] = "post"
        results = append(results, result)
    }

    // Search Projects
    projectsCursor, err := projectsCollection.Find(ctx, filter, findOptions)
    if err != nil {
        log.Printf("Error searching projects: %v", err)
        return nil, 0, err
    }
    defer projectsCursor.Close(ctx)

    for projectsCursor.Next(ctx) {
        var result bson.M
        if err := projectsCursor.Decode(&result); err != nil {
            log.Printf("Error decoding project result: %v", err)
            return nil, 0, err
        }
        result["type"] = "project"
        results = append(results, result)
    }

    // Count total results
    postsCount, err := postsCollection.CountDocuments(ctx, filter)
    if err != nil {
        log.Printf("Error counting posts: %v", err)
        return nil, 0, err
    }

    projectsCount, err := projectsCollection.CountDocuments(ctx, filter)
    if err != nil {
        log.Printf("Error counting projects: %v", err)
        return nil, 0, err
    }

    totalCount = postsCount + projectsCount

    log.Printf("Found %d total results", totalCount)

    return results, totalCount, nil
}
