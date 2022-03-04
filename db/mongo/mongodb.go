package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Queries struct {
	users    *mongo.Collection
	movies   *mongo.Collection
	comments *mongo.Collection
	sessions *mongo.Collection
	theaters *mongo.Collection
}

func NewMongoQueries(db *mongo.Database) *Queries {
	// Require that the data inserted into the 'users' collection
	// must contain both 'email' and 'password' fields
	usersValidatorModels := &options.CreateCollectionOptions{Validator: bson.D{{
		"$jsonSchema", bson.D{
			{"bsonType", "object"},
			{"required", []string{"email", "password"}},
			{"properties", bson.D{
				{"email", bson.D{
					{"bsonType", "string"},
					{"description", "must be a string and is required"},
				}},
				{"password", bson.D{
					{"bsonType", "string"},
					{"description", "must be a string and is required"},
				}},
			}},
		},
	}}}
	createSchemaValidation(db, "users", usersValidatorModels)

	// Create unique indexes for the 'name' and 'email' fields
	// in the 'users' collection
	userIndexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	AddIndexMany(db, "users", userIndexModels)

	// Requires that the data inserted into the 'movies' collection
	// must contain the 'title' field
	moviesValidatorModels := &options.CreateCollectionOptions{Validator: bson.D{{
		"$jsonSchema", bson.D{
			{"bsonType", "object"},
			{"required", []string{"title"}},
			{"properties", bson.D{{
				"title", bson.D{
					{"bsonType", "string"},
					{"description", "must be a string and is required"},
				},
			}}},
		},
	}}}

	// Create text indexes for the 'cast', 'fullplot', 'genres',
	// and 'title' fields in the 'movies' collection,
	// as well as a general index for the runtime field
	createSchemaValidation(db, "movies", moviesValidatorModels)
	moviesIndexModels := []mongo.IndexModel{
		{Keys: bson.M{"runtime": -1}},
		{Keys: bson.D{
			{"cast", "text"},
			{"fullplot", "text"},
			{"genres", "text"},
			{"title", "text"},
		}},
	}
	AddIndexMany(db, "movies", moviesIndexModels)

	// Requires that the data inserted into the 'comments' collection
	// must contain the 'text' field
	commentsValidatorModels := &options.CreateCollectionOptions{Validator: bson.D{{
		"$jsonSchema", bson.D{
			{"bsonType", "object"},
			{"required", []string{"text"}},
			{"properties", bson.D{{
				"text", bson.D{
					{"bsonType", "string"},
					{"description", "must be a string and is required"},
				},
			}}},
		},
	}}}
	createSchemaValidation(db, "comments", commentsValidatorModels)

	return &Queries{
		users:    db.Collection("users"),
		movies:   db.Collection("movies"),
		comments: db.Collection("comments"),
		sessions: db.Collection("sessions"),
		theaters: db.Collection("theaters"),
	}
}

// createSchemaValidation Create schemaValidation for collections
func createSchemaValidation(db *mongo.Database, collectionName string, validators *options.CreateCollectionOptions) {
	db.CreateCollection(context.Background(), collectionName, validators)
}

// AddIndexOne Creating indexes for fields in a collection
func AddIndexOne(db *mongo.Database, collectionName string, indexKeys mongo.IndexModel) {
	db.Collection(collectionName).Indexes().CreateOne(context.Background(), indexKeys)

}

// AddIndexMany Create multiple indexes for fields in a collection
func AddIndexMany(db *mongo.Database, collectionName string, indexKeys []mongo.IndexModel) {
	db.Collection(collectionName).Indexes().CreateMany(context.Background(), indexKeys)
}
