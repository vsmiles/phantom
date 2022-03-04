package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"phantom/util"
	"testing"
	"time"
)

var testQueries *Queries
var testDatabase *mongo.Database

func TestMain(m *testing.M) {
	// Load Config
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// Connect mongodb, Timeout: 10ms
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoSource).SetServerAPIOptions(serverAPIOptions))
	defer cancel()
	if err != nil {
		log.Fatal("cannot connect mongodb: ", err)
	}

	testDatabase = mongoClient.Database("phantom")
	testQueries = NewMongoQueries(testDatabase)

	os.Exit(m.Run())
}
