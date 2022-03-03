package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"phantom/util"
	"testing"
)

func randomUser() AddUserParams {
	return AddUserParams{
		Name:     util.RandomUser(),
		Email:    util.RandomEmail(),
		Password: util.RandomString(6),
	}
}

func addUser(t *testing.T, user AddUserParams) primitive.ObjectID {
	id, err := testQueries.AddUser(context.Background(), user)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	return id
}

func TestAddUser(t *testing.T) {
	user1 := randomUser()
	addUser(t, user1)
	user2 := randomUser()
	user2.Name = user1.Name
	_, err := testQueries.AddUser(context.Background(), user2)
	require.Error(t, err)
	user3 := randomUser()
	user3.Email = user1.Email
	_, err = testQueries.AddUser(context.Background(), user3)
	require.Error(t, err)
}

func getUserByID(t *testing.T, id primitive.ObjectID) User {
	user, err := testQueries.GetUserByID(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

func TestGetUserByID(t *testing.T) {
	user1 := randomUser()
	id := addUser(t, user1)
	user2 := getUserByID(t, id)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
}

func TestGetUserByName(t *testing.T) {
	user1 := randomUser()
	id := addUser(t, user1)
	user2, err := testQueries.GetUserByName(context.Background(), user1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, id, user2.ID)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	_, err = testQueries.GetUserByName(context.Background(), util.RandomUser())
	require.Error(t, err)
	require.EqualError(t, mongo.ErrNoDocuments, err.Error())
}

func TestGetUserByEmail(t *testing.T) {
	user1 := randomUser()
	id := addUser(t, user1)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, id, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Password, user2.Password)
}

func TestUpdateUserName(t *testing.T) {
	id := addUser(t, randomUser())
	user1 := getUserByID(t, id)
	user1.Name = util.RandomUser()
	updateResult, err := testQueries.UpdateUserName(context.Background(), user1)
	require.NoError(t, err)
	require.NotEmpty(t, updateResult)
	user2 := getUserByID(t, id)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Password, user2.Password)
}

func TestUpdateUserPassword(t *testing.T) {
	id := addUser(t, randomUser())
	user1 := getUserByID(t, id)
	user1.Password = util.RandomString(6)
	updateResult, err := testQueries.UpdateUserPassword(context.Background(), user1)
	require.NoError(t, err)
	require.NotEmpty(t, updateResult)
	user2 := getUserByID(t, id)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Password, user2.Password)
}
