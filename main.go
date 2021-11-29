// Recipes API
//
// Это небольшой проект на Go с рецептами.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Алексей Кузьмин <kuzmin.al35@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"context"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"recipes-api/handlers"
)

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection
var mongoDbUri = os.Getenv("MONGO_URI")
var mongoDbDatabase = os.Getenv("MONGODB_DATABASE")
var mongoDbCollection = os.Getenv("MONGODB_COLLECTION")

var recipesHandler *handlers.RecipesHandler

var authHandler *handlers.AuthHandler

func init() {
	os.Setenv("JWT_SECRET", "eUbP9shywUygMx7u")
	if mongoDbUri == "" {
		mongoDbUri = "mongodb://localhost:27017/"
	}
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDbUri))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	if mongoDbDatabase == "" {
		mongoDbDatabase = "recipe"
	}
	if mongoDbCollection == "" {
		mongoDbCollection = "recipes"
	}
	collection = client.Database(mongoDbDatabase).Collection(mongoDbCollection)
	log.Println("Connected to MongoDB")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	collectionUsers := client.Database(mongoDbDatabase).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/users", authHandler.CreateUser)

	auth := router.Group("/")
	auth.Use(authHandler.AuthMiddleware())
	{
		auth.POST("/recipes", recipesHandler.NewRecipeHandler)
		auth.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
		//router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
		auth.PUT("/recipes:id", recipesHandler.UpdateRecipeHandler)
		auth.DELETE("/recipes:id", recipesHandler.DeleteRecipeHandler)
	}

	router.Run()
}
