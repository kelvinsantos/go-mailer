package services

import "go.mongodb.org/mongo-driver/mongo"

var Client *mongo.Client

func SetClient(client *mongo.Client) {
	Client = client
}
