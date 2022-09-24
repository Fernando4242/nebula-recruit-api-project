package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client = Connection()

func main() {
	router := gin.Default()
	// main routes
	router.GET("/create/:username/:password", createAccount)
	router.GET("/read", getAllAccounts)
	router.GET("/update/:id/:username", updateAccount)
	router.GET("/delete/:id", deleteAccount)
	router.Run("localhost:3000")
}

// creates an account with the parameters passed into the url
// saves it in a mongodb repo
func createAccount(c *gin.Context) {
	collection := client.Database("test").Collection("testing")
	username := c.Param("username")
	password := c.Param("password")

	if username == "" || password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing username or password"})
		return
	}

	result, err := collection.InsertOne(
		context.TODO(),
		bson.D{
			{"username", username},
			{"password", password}},
	)

	if err != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Error while uploading document"})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// gets all the accounts saved on the database
func getAllAccounts(c *gin.Context) {
	collection := client.Database("test").Collection("testing")

	filter := bson.D{}
	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Error while retrieving all documents"})
		return
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Failed to clean up data"})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

// updates an accounts username based on id
func updateAccount(c *gin.Context) {
	collection := client.Database("test").Collection("testing")

	id := c.Param("id")
	username := c.Param("username")

	if id == "" || username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id or username"})
		return
	}

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", _id}}
	update := bson.D{{"$set", bson.D{{"username", username}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Error while uploading document"})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// deletes an account based on its id
func deleteAccount(c *gin.Context) {
	collection := client.Database("test").Collection("testing")

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id"})
		return
	}

	_id, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.DeleteOne(context.TODO(), bson.D{{"_id", _id}})
	if err != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Error while uploading document"})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// creates the instance of the database
// I wasn't sure how so I did it this way
// I always struggle with the best way to do this
func Connection() *mongo.Client {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("No .en file found")
	}

	uri := os.Getenv("DB_URI")
	if uri == "" {
		log.Fatalln("Missing DB URI")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to database succesfully")
	return client
}
