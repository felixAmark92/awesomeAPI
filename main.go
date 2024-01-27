package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"

// album represents data about a record album.
type album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var client *mongo.Client

func main() {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(),
		bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	println("Pinged your deployment. You successfully connected to MongoDB!")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	err = router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {

	coll := client.Database("goDb").Collection("albums")

	filter := bson.D{{}}

	result, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	defer result.Close(context.TODO())

	var albums []bson.M
	for result.Next(context.TODO()) {
		var album bson.M
		if err := result.Decode(&album); err != nil {
			panic(err)
		}
		albums = append(albums, album)
	}

	if err := result.Err(); err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	//establish connection to your MongoDb server
	coll := client.Database("goDb").Collection("albums")
	// Add the new album to the mongoDb.
	result, err := coll.InsertOne(context.TODO(), newAlbum)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, result)
}
