package db

import (
	"context"
	"log"

	"github.com/lscantillo/twitter-clone-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCN = ConnectMongoDB()
var clientOptions = options.Client().ApplyURI(config.GetVariables("MONGO_URL"))

func ConnectMongoDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

// Check the connection
func CheckDB() int {
	err := MongoCN.Ping(context.TODO(), nil)

	if err != nil {
		return 0
	}
	return 1

}
