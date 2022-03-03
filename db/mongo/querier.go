package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Querier interface {
	AddUser(ctx context.Context, arg AddUserParams) (primitive.ObjectID, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (User, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUserName(ctx context.Context, user User) (*mongo.UpdateResult, error)
	UpdateUserPassword(ctx context.Context, user User) (*mongo.UpdateResult, error)
	AddComment(ctx context.Context, arg AddCommentParams) (primitive.ObjectID, error)
	GetComment(ctx context.Context, id primitive.ObjectID) (Comments, error)
	GetCommentsByMovieID(ctx context.Context, arg GetCommentsParams) ([]Comments, error)
	GetCommentsByName(ctx context.Context, arg GetCommentsParams) ([]Comments, error)
	UpdateComment(ctx context.Context, comment Comments) (*mongo.UpdateResult, error)
	DeleteComment(ctx context.Context, id primitive.ObjectID, name string) (int64, error)
	AddMovie(ctx context.Context, arg AddMovieParams) (primitive.ObjectID, error)
	AddMovies(ctx context.Context, arg []AddMovieParams) ([]interface{}, error)
	GetMovieByID(ctx context.Context, id primitive.ObjectID) (Movies, error)
	GetMoviesByGenres(ctx context.Context, arg GetMoviesParams) ([]Movies, error)
	SearchForMovies(ctx context.Context, arg SearchForMoviesParams) ([]Movies, error)
	GetTheMostViewedMovies(ctx context.Context, arg GetMoviesParams) ([]Movies, error)
	GetTheLatestReleasedMovies(ctx context.Context, arg GetMoviesParams) ([]Movies, error)
	ReplaceMovieInfoByID(ctx context.Context, id primitive.ObjectID, movie Movies) (*mongo.UpdateResult, error)
}

var _ Querier = (*Queries)(nil)
