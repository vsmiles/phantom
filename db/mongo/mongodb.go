package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Queries struct {
	users    *mongo.Collection
	movies   *mongo.Collection
	comments *mongo.Collection
	sessions *mongo.Collection
	theaters *mongo.Collection
}

func NewMongoQueries(db *mongo.Database) *Queries {
	return &Queries{
		users:    db.Collection("users"),
		movies:   db.Collection("movies"),
		comments: db.Collection("comments"),
		sessions: db.Collection("sessions"),
		theaters: db.Collection("theaters"),
	}
}
