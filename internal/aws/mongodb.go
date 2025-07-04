package aws

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup
var MongoDBClient *mongo.Client

type ColIndex struct {
	cn      string             // Collection name
	indexes []mongo.IndexModel // Collection list of indexes
}

func (r ColIndex) CreateCollectionIndexes(db *mongo.Database, err chan error) {
	defer close(err)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	col := db.Collection(r.cn)
	if _, e := col.Indexes().CreateMany(ctx, r.indexes); e != nil {
		err <- e
		return
	}

	err <- nil
}

func connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags:   true, // Replaces the bson struct tags with json struct tags
		ObjectIDAsHexString: true, // Allows the ObjectID to be marshalled as a string
		NilSliceAsEmpty:     true,
		UseLocalTimeZone:    false,
	}

	client, err := mongo.Connect(options.Client().SetBSONOptions(bsonOpts).ApplyURI(os.Getenv("AWS_MONGO_DB_URL")))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	MongoDBClient = client

	// Creating Indexes for all collections in th db
	var allPossibleCollectionsIndexes []ColIndex
	allPossibleCollectionsIndexes = []ColIndex{
		{
			cn: "roles",
			indexes: []mongo.IndexModel{
				{
					Keys: bson.D{{"name", 1}},
				},
				{
					Keys: bson.D{{"created_at", -1}},
				},
				{
					Keys: bson.D{{"updated_at", -1}},
				},
			},
		},
		{
			cn: "teams",
			indexes: []mongo.IndexModel{
				{
					Keys: bson.D{{"name", 1}},
				},
				{
					Keys: bson.D{{"name", "text"}, {"description", "text"}},
				},
				{
					Keys: bson.D{{"created_at", -1}},
				},
				{
					Keys: bson.D{{"updated_at", -1}},
				},
			},
		},
		{
			cn: "users",
			indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"email", 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{"first_name", "text"}, {"last_name", "text"}, {"email", "text"}, {"full_name", "text"}},
				},
			},
		},
	}

	errChan := make(chan error, len(allPossibleCollectionsIndexes))

	// Perform a parallel computation for creating all the possible index for our collection
	for _, allPossibleCollectionIndexes := range allPossibleCollectionsIndexes {
		wg.Add(1)

		go func(i ColIndex) {
			defer func() {
				recover()
			}()

			defer wg.Done()

			i.CreateCollectionIndexes(client.Database("flowCx"), errChan)
			<-errChan

			close(errChan)

		}(allPossibleCollectionIndexes)
	}

	wg.Wait()
	return client, nil
}

func ConnectMongoDB() (*mongo.Client, error) {
	client, err := connect()
	MongoDBClient = client

	return client, err
}
