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
	"github.com/gin-gonic/gin"
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

func init() {
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
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
	//router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.PUT("/recipes:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes:id", recipesHandler.DeleteRecipeHandler)
	router.Run()
}
