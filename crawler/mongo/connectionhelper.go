package mongo

import (
	"context"
	"entities"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

const (
	CONNECTIONSTRING = "mongodb://mongo:27017"
)

//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func CreateDocument(client *mongo.Client, posto entities.Pessoajuridica) error {
	db := os.Getenv("MONGO_DATABASE")
	if db == "" {
		log.Fatal("É necessário configurar a variável de ambiente MONGO_DATABASE")
	}
	collectionName := os.Getenv("MONGO_COLLECTION")
	if collectionName == "" {
		log.Fatal("É necessário configurar a variável de ambiente MONGO_COLLECTION.")
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(db).Collection(collectionName)
	//Perform InsertOne operation & validate against the error.
	_, err := collection.InsertOne(context.TODO(), posto)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}
