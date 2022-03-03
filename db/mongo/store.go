package db

import "go.mongodb.org/mongo-driver/mongo"

type Store interface {
	Querier
}

type MongoStore struct {
	*Queries
}

func NewMongoStore(database *mongo.Database) *MongoStore {
	return &MongoStore{
		Queries: NewMongoQueries(database),
	}
}
