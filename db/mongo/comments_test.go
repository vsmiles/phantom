package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"phantom/util"
	"testing"
)

func randomComment(user User, movie Movies) AddCommentParams {
	return AddCommentParams{
		Name:    user.Name,
		Email:   user.Email,
		MovieID: movie.Id,
		Text:    util.RandomString(140),
	}
}

func addComment(t *testing.T, user User, movie Movies) primitive.ObjectID {
	id, err := testQueries.AddComment(context.Background(), randomComment(user, movie))
	require.NoError(t, err)
	require.NotEmpty(t, id)
	return id
}

func TestAddComment(t *testing.T) {
	user := getUserByID(t, addUser(t, randomUser()))
	movie := getMovieByID(t, addMovie(t, randomMovie()))
	addComment(t, user, movie)
}

func getCommentByID(t *testing.T, id primitive.ObjectID) Comments {
	comment, err := testQueries.GetComment(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	return comment
}

func TestGetComment(t *testing.T) {
	user := getUserByID(t, addUser(t, randomUser()))
	movie := getMovieByID(t, addMovie(t, randomMovie()))
	commentId := addComment(t, user, movie)
	getCommentByID(t, commentId)
}

func TestGetCommentsByMovieID(t *testing.T) {
	n := 10
	movie := getMovieByID(t, addMovie(t, randomMovie()))
	for i := 0; i < n; i++ {
		user := getUserByID(t, addUser(t, randomUser()))
		addComment(t, user, movie)
	}
	arg := GetCommentsParams{
		MovieID: movie.Id,
		Limit:   5,
		Skip:    5,
	}
	comments, err := testQueries.GetCommentsByMovieID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Equal(t, 5, len(comments))
	for i := range comments {
		require.NotEmpty(t, comments[i])
		require.Equal(t, movie.Id, comments[i].MovieID)
	}
}

func TestGetCommentsByName(t *testing.T) {
	n := 10
	user := getUserByID(t, addUser(t, randomUser()))
	for i := 0; i < n; i++ {
		movie := getMovieByID(t, addMovie(t, randomMovie()))
		addComment(t, user, movie)
	}
	arg := GetCommentsParams{
		Name:  user.Name,
		Limit: 5,
		Skip:  5,
	}
	comments, err := testQueries.GetCommentsByName(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Equal(t, 5, len(comments))
	for i := range comments {
		require.NotEmpty(t, comments[i])
		require.Equal(t, user.Name, comments[i].Name)
	}
}

func TestUpdateComment(t *testing.T) {
	user1 := getUserByID(t, addUser(t, randomUser()))
	movie1 := getMovieByID(t, addMovie(t, randomMovie()))
	comment1 := getCommentByID(t, addComment(t, user1, movie1))
	comment1.Text = util.RandomString(140)
	updateResult, err := testQueries.UpdateComment(context.Background(), comment1)
	require.NoError(t, err)
	require.NotEmpty(t, updateResult)
	comment2 := getCommentByID(t, comment1.ID)
	require.Equal(t, comment1.Text, comment2.Text)

	user2 := getUserByID(t, addUser(t, randomUser()))
	movie2 := getMovieByID(t, addMovie(t, randomMovie()))
	comment3 := getCommentByID(t, addComment(t, user2, movie2))
	comment4 := comment3

	comment3.ID = primitive.NewObjectID()
	_, err = testQueries.UpdateComment(context.Background(), comment3)
	require.Equal(t, mongo.ErrNoDocuments, err)

	comment4.Name = util.RandomString(6)
	_, err = testQueries.UpdateComment(context.Background(), comment4)
	require.Equal(t, mongo.ErrNoDocuments, err)
}

func TestDeleteComment(t *testing.T) {
	comment1 := getCommentByID(t, addComment(t, getUserByID(t, addUser(t, randomUser())), getMovieByID(t, addMovie(t, randomMovie()))))
	deletedCount, err := testQueries.DeleteComment(context.Background(), comment1.ID, comment1.Name)
	require.NoError(t, err)
	require.Equal(t, int64(1), deletedCount)
	comment2 := getCommentByID(t, addComment(t, getUserByID(t, addUser(t, randomUser())), getMovieByID(t, addMovie(t, randomMovie()))))
	comment3 := comment2
	comment2.ID = primitive.NewObjectID()
	_, err = testQueries.DeleteComment(context.Background(), comment1.ID, comment1.Name)
	require.Equal(t, mongo.ErrNoDocuments, err)

	comment3.Name = util.RandomString(6)
	_, err = testQueries.DeleteComment(context.Background(), comment1.ID, comment1.Name)
	require.Equal(t, mongo.ErrNoDocuments, err)
}
