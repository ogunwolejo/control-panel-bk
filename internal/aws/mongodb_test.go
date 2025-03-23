package aws

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoTestSuite defines a test suite for MongoDB
type MongoTestSuite struct {
	suite.Suite
	client *mongo.Client
	db     *mongo.Database
}

// SetupSuite runs before all test cases
func (suite *MongoTestSuite) SetupSuite() {
	// Set a test MongoDB URI (use a real test URI or a mocked one)
	testMongoURI := os.Getenv("AWS_MONGO_DB_URL_TEST")
	if testMongoURI == "" {
		testMongoURI = "mongodb://localhost:27017"
	}

	// Connect to MongoDB
	_, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	suite.Require().NoError(err, "MongoDB connection failed")

	suite.client = client
	suite.db = client.Database("testDb")
}

// TearDownSuite runs after all test cases
func (suite *MongoTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	if suite.client != nil {
		err := suite.client.Disconnect(ctx)
		suite.Require().NoError(err, "Failed to disconnect MongoDB client")
	}
}

// TestConnectMongoDB checks if `ConnectMongoDB` works correctly
func (suite *MongoTestSuite) TestConnectMongoDB() {
	client, err := ConnectMongoDB()
	assert.NoError(suite.T(), err, "ConnectMongoDB should not return an error")
	assert.NotNil(suite.T(), client, "MongoDB client should not be nil")
}

// TestCreateCollectionIndexes checks if indexes are created correctly
func (suite *MongoTestSuite) TestCreateCollectionIndexes() {
	collectionName := "test_collection"

	// Define mock indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"name", 1}},
		},
		{
			Keys: bson.D{{"created_at", -1}},
		},
	}

	colIndex := ColIndex{
		cn:      collectionName,
		indexes: indexes,
	}

	errChan := make(chan error)
	go colIndex.CreateCollectionIndexes(suite.db, errChan)

	err := <-errChan
	assert.NoError(suite.T(), err, "Index creation should not return an error")

	// Verify indexes exist
	cursor, err := suite.db.Collection(collectionName).Indexes().List(context.TODO())
	assert.NoError(suite.T(), err, "Should be able to list indexes")

	defer cursor.Close(context.TODO())

	var foundIndexes []bson.M
	err = cursor.All(context.TODO(), &foundIndexes)
	assert.NoError(suite.T(), err, "Should decode index list successfully")
	assert.GreaterOrEqual(suite.T(), len(foundIndexes), 2, "Should have at least 2 indexes")
}

// Run the test suite
func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
