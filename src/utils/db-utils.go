package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"kelvin.com/mailer/src/env"
)

func GetClient() *mongo.Client {
	log.Println("Opening connection to database...")
	client_options := options.Client().ApplyURI(env.GO_MAILER_DB_URI)
	client, err := mongo.NewClient(client_options)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("Cannot connect to database!", err)
		return nil
	} else {
		log.Println("Succesfully connected to database!")
		return client
	}
}
