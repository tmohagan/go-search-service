package db

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

func VerifyDatabaseContent() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    database := Client.Database("test") // Replace with your actual database name
    postsCollection := database.Collection("posts")
    projectsCollection := database.Collection("projects")

    postCount, err := postsCollection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return fmt.Errorf("error counting posts: %v", err)
    }

    projectCount, err := projectsCollection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return fmt.Errorf("error counting projects: %v", err)
    }

    log.Printf("Database contains %d posts and %d projects", postCount, projectCount)

    // Print a sample document from posts collection
    var samplePost bson.M
    err = postsCollection.FindOne(ctx, bson.M{}).Decode(&samplePost)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Println("No posts found in the database")
        } else {
            return fmt.Errorf("error retrieving sample post: %v", err)
        }
    } else {
        log.Printf("Sample post: %+v", samplePost)
    }

    return nil
}

func PerformSearch(query string) ([]bson.M, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    database := Client.Database("your_database_name") // Replace with your actual database name
    postsCollection := database.Collection("posts")
    projectsCollection := database.Collection("projects")

    log.Printf("Searching for query: %s", query)

    filter := bson.M{
        "$or": []bson.M{
            {"title": bson.M{"$regex": ".*" + query + ".*", "$options": "i"}},
            {"summary": bson.M{"$regex": ".*" + query + ".*", "$options": "i"}},
            {"content": bson.M{"$regex": ".*" + query + ".*", "$options": "i"}},
        },
    }

    log.Printf("Search filter: %+v", filter)

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
        log.Printf("Found post: %+v", result)
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
        log.Printf("Found project: %+v", result)
        results = append(results, result)
    }

    log.Printf("Found %d total results", len(results))

    return results, nil
}