package main

import (
    "context"
    "fmt"
    "log"
    "mongodb-api-demo/controllers"
    "mongodb-api-demo/services"
    "os"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
    server         *gin.Engine
    userService    services.UserService
    userController controllers.UserController
    ctx            context.Context
    userCollection *mongo.Collection
    mongoClient    *mongo.Client
    err            error
)

func init() {
    ctx := context.TODO()
    uri := fmt.Sprintf(
        "mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
        os.Getenv("mongodb_username"),
        os.Getenv("mongodb_password"),
        os.Getenv("mongodb_host"),
    )
    mongoConnection := options.Client().ApplyURI(uri)
    mongoClient, err := mongo.Connect(ctx, mongoConnection)
    if err != nil {
        log.Fatal(err)
    }

    err = mongoClient.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("mongo connection established")
    userCollection = mongoClient.Database("userdb").Collection("users")
    userService = services.NewUserService(userCollection, ctx)
    userController = controllers.New(userService)
    server = gin.Default()
}

func main() {
    defer mongoClient.Disconnect(ctx)

    basePath := server.Group("/v1")
    userController.RegisterUserRoutes(basePath)

    log.Fatal(server.Run(":9090"))
}
