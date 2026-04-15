package helpers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("ecommerce").Collection(collectionName)
	return collection
}
