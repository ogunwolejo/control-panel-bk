package aws

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"os"
	"time"
)

var MongoDBClient *mongo.Client
var MongoFlowCxDBClient *mongo.Database

func connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("AWS_MONGO_DB_URL")))
	if err != nil {
		return nil, err
	}

	client.Ping(ctx, readpref.Primary())

	MongoFlowCxDBClient = client.Database("flowCx")
	return client, nil
}

func ConnectMongoDB() (*mongo.Client, error) {
	client, err := connect()
	MongoDBClient = client

	return client, err
}
