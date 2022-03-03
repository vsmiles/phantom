package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Comments struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Email   string             `json:"email" bson:"email"`
	MovieID primitive.ObjectID `json:"movie_id" bson:"movie_id"`
	Text    string             `json:"text" bson:"text"`
	Date    primitive.DateTime `json:"date" bson:"date"`
}

type AddCommentParams struct {
	Name    string             `json:"name" bson:"name,omitempty"`
	Email   string             `json:"email" bson:"email,omitempty"`
	MovieID primitive.ObjectID `json:"movie_id" bson:"movie_id,omitempty"`
	Text    string             `json:"text" bson:"text,omitempty"`
}

type GetCommentsParams struct {
	Name    string             `json:"name" bson:"name"`
	MovieID primitive.ObjectID `json:"movie_id" bson:"movie_id,omitempty"`
	Limit   int64              `json:"limit"`
	Skip    int64              `json:"skip"`
}

func (q *Queries) AddComment(ctx context.Context, arg AddCommentParams) (primitive.ObjectID, error) {
	res, err := q.comments.InsertOne(ctx, arg)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (q Queries) GetComment(ctx context.Context, id primitive.ObjectID) (Comments, error) {
	var comment Comments
	err := q.comments.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return Comments{}, err
	}
	return comment, nil
}

func (q Queries) GetCommentsByMovieID(ctx context.Context, arg GetCommentsParams) ([]Comments, error) {
	var findOptions *options.FindOptions
	if arg.Limit > 0 {
		findOptions = &options.FindOptions{}
		findOptions.SetLimit(arg.Limit)
		findOptions.SetSkip(arg.Skip)
	}
	cursor, err := q.comments.Find(ctx, bson.M{"movie_id": arg.MovieID}, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var comments []Comments
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (q Queries) GetCommentsByName(ctx context.Context, arg GetCommentsParams) ([]Comments, error) {
	var findOptions *options.FindOptions
	if arg.Limit > 0 {
		findOptions = &options.FindOptions{}
		findOptions.SetLimit(arg.Limit)
		findOptions.SetSkip(arg.Skip)
	}
	cursor, err := q.comments.Find(ctx, bson.M{"name": arg.Name}, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var comments []Comments
	if err = cursor.All(ctx, &comments); err != nil {
		return []Comments{}, err
	}
	return comments, nil
}

func (q Queries) UpdateComment(ctx context.Context, comment Comments) (*mongo.UpdateResult, error) {
	res, err := q.comments.UpdateOne(ctx,
		bson.D{
			{"_id", comment.ID},
			{"name", comment.Name},
		},
		bson.D{
			{
				"$set", bson.D{
					{"text", comment.Text},
				},
			},
		})
	if err != nil {
		return nil, err
	}
	if res.ModifiedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return res, nil
}

func (q *Queries) DeleteComment(ctx context.Context, id primitive.ObjectID, name string) (int64, error) {
	deleteResult, err := q.comments.DeleteOne(ctx, bson.D{{"_id", id}, {"name", name}})
	if err != nil {
		return 0, err
	}
	if deleteResult.DeletedCount == 0 {
		return 0, mongo.ErrNoDocuments
	}
	return deleteResult.DeletedCount, nil

}
