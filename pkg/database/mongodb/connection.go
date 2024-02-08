package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setup the mongodb connection
type MongoDBConn struct {
	client *mongo.Client
}

// the function to create a new mongodb connection
func NewMongoDBConn(uri string) (*MongoDBConn, error) {
	// use the URI to connect to the database
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoDBConn{client: client}, nil
}
