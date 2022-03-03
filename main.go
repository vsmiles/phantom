package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"phantom/api"
	db "phantom/db/mongo"
	"phantom/util"
	"time"
)

func main() {
	// Load Config
	config, err := util.LoadConfig(".")
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

	mongoDatabase := mongoClient.Database("sample_mflix")
	mongoStore := db.NewMongoStore(mongoDatabase)
	server, err := api.NewServer(config, mongoStore)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
