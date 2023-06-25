package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDatabase() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(Env("MONGOURI")))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	//ping database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

var Database *mongo.Client = ConnectDatabase()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("poosible_db").Collection(collectionName)
}
