package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
}

type AddUserParams struct {
	Name     string `json:"name" bson:"name,omitempty"`
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) (primitive.ObjectID, error) {
	res, err := q.users.InsertOne(ctx, arg)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (q *Queries) GetUserByID(ctx context.Context, id primitive.ObjectID) (User, error) {
	var user User
	err := q.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (q *Queries) GetUserByName(ctx context.Context, name string) (User, error) {
	var user User
	err := q.users.FindOne(ctx, bson.M{"name": name}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := q.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (q *Queries) UpdateUserName(ctx context.Context, user User) (*mongo.UpdateResult, error) {

	res, err := q.users.UpdateByID(ctx, user.ID, bson.D{
		{"$set", bson.D{
			{"name", user.Name},
		}},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (q *Queries) UpdateUserPassword(ctx context.Context, user User) (*mongo.UpdateResult, error) {

	res, err := q.users.UpdateByID(ctx, user.ID, bson.D{
		{"$set", bson.D{
			{"password", user.Password},
		}},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
