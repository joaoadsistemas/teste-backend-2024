package db

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection
var client *mongo.Client
var once sync.Once

func Connection() *mongo.Collection {
	var err error
	once.Do(func() {
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			fmt.Println("MONGO: ", err)
			return
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			fmt.Println("MONGO: ", err)
			return
		}

		db = client.Database("teste_backend").Collection("products")
	})

	return db
}

func Disconnect() {
	if client != nil {
		client.Disconnect(context.TODO())
	}
}