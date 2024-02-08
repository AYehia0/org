package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setup the mongodb connection
type MongoDBConn struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// the function to create a new mongodb connection
func NewMongoDBConn(uri, databaseName, username, password string) (*MongoDBConn, error) {

	// authenticate the connection
	credentials := options.Credential{
		Username: username,
		Password: password,
	}

	// use the URI to connect to the database
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri).SetAuth(credentials))
	if err != nil {
		return nil, err
	}

	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// access the database
	database := client.Database(databaseName)

	return &MongoDBConn{Client: client, Database: database}, nil
}
