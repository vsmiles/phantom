package db

import "go.mongodb.org/mongo-driver/mongo"

type Queries struct {
	database *mongo.Database
	users    *mongo.Collection
	movies   *mongo.Collection
	comments *mongo.Collection
	sessions *mongo.Collection
	theaters *mongo.Collection
}

func NewMongoQueries(database *mongo.Database) *Queries {
	return &Queries{
		database: database,
		users:    database.Collection("users"),
		movies:   database.Collection("movies"),
		comments: database.Collection("comments"),
		sessions: database.Collection("sessions"),
		theaters: database.Collection("theaters"),
	}
}
