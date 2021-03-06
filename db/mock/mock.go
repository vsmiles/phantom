// Code generated by MockGen. DO NOT EDIT.
// Source: phantom/db/mongo (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	mongo0 "phantom/db/mongo"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddComment mocks base method.
func (m *MockStore) AddComment(arg0 context.Context, arg1 mongo0.AddCommentParams) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", arg0, arg1)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddComment indicates an expected call of AddComment.
func (mr *MockStoreMockRecorder) AddComment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockStore)(nil).AddComment), arg0, arg1)
}

// AddMovie mocks base method.
func (m *MockStore) AddMovie(arg0 context.Context, arg1 mongo0.AddMovieParams) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMovie", arg0, arg1)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMovie indicates an expected call of AddMovie.
func (mr *MockStoreMockRecorder) AddMovie(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMovie", reflect.TypeOf((*MockStore)(nil).AddMovie), arg0, arg1)
}

// AddMovies mocks base method.
func (m *MockStore) AddMovies(arg0 context.Context, arg1 []mongo0.AddMovieParams) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMovies", arg0, arg1)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMovies indicates an expected call of AddMovies.
func (mr *MockStoreMockRecorder) AddMovies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMovies", reflect.TypeOf((*MockStore)(nil).AddMovies), arg0, arg1)
}

// AddUser mocks base method.
func (m *MockStore) AddUser(arg0 context.Context, arg1 mongo0.AddUserParams) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStoreMockRecorder) AddUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockStore)(nil).AddUser), arg0, arg1)
}

// DeleteComment mocks base method.
func (m *MockStore) DeleteComment(arg0 context.Context, arg1 primitive.ObjectID, arg2 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComment", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockStoreMockRecorder) DeleteComment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockStore)(nil).DeleteComment), arg0, arg1, arg2)
}

// GetComment mocks base method.
func (m *MockStore) GetComment(arg0 context.Context, arg1 primitive.ObjectID) (mongo0.Comments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComment", arg0, arg1)
	ret0, _ := ret[0].(mongo0.Comments)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComment indicates an expected call of GetComment.
func (mr *MockStoreMockRecorder) GetComment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComment", reflect.TypeOf((*MockStore)(nil).GetComment), arg0, arg1)
}

// GetCommentsByMovieID mocks base method.
func (m *MockStore) GetCommentsByMovieID(arg0 context.Context, arg1 mongo0.GetCommentsParams) ([]mongo0.Comments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentsByMovieID", arg0, arg1)
	ret0, _ := ret[0].([]mongo0.Comments)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentsByMovieID indicates an expected call of GetCommentsByMovieID.
func (mr *MockStoreMockRecorder) GetCommentsByMovieID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentsByMovieID", reflect.TypeOf((*MockStore)(nil).GetCommentsByMovieID), arg0, arg1)
}

// GetCommentsByName mocks base method.
func (m *MockStore) GetCommentsByName(arg0 context.Context, arg1 mongo0.GetCommentsParams) ([]mongo0.Comments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentsByName", arg0, arg1)
	ret0, _ := ret[0].([]mongo0.Comments)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentsByName indicates an expected call of GetCommentsByName.
func (mr *MockStoreMockRecorder) GetCommentsByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentsByName", reflect.TypeOf((*MockStore)(nil).GetCommentsByName), arg0, arg1)
}

// GetMovieByID mocks base method.
func (m *MockStore) GetMovieByID(arg0 context.Context, arg1 primitive.ObjectID) (mongo0.Movies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovieByID", arg0, arg1)
	ret0, _ := ret[0].(mongo0.Movies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovieByID indicates an expected call of GetMovieByID.
func (mr *MockStoreMockRecorder) GetMovieByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovieByID", reflect.TypeOf((*MockStore)(nil).GetMovieByID), arg0, arg1)
}

// GetMoviesByGenres mocks base method.
func (m *MockStore) GetMoviesByGenres(arg0 context.Context, arg1 mongo0.GetMoviesParams) ([]mongo0.Movies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesByGenres", arg0, arg1)
	ret0, _ := ret[0].([]mongo0.Movies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesByGenres indicates an expected call of GetMoviesByGenres.
func (mr *MockStoreMockRecorder) GetMoviesByGenres(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesByGenres", reflect.TypeOf((*MockStore)(nil).GetMoviesByGenres), arg0, arg1)
}

// GetTheLatestReleasedMovies mocks base method.
func (m *MockStore) GetTheLatestReleasedMovies(arg0 context.Context, arg1 mongo0.GetMoviesParams) ([]mongo0.Movies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTheLatestReleasedMovies", arg0, arg1)
	ret0, _ := ret[0].([]mongo0.Movies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTheLatestReleasedMovies indicates an expected call of GetTheLatestReleasedMovies.
func (mr *MockStoreMockRecorder) GetTheLatestReleasedMovies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTheLatestReleasedMovies", reflect.TypeOf((*MockStore)(nil).GetTheLatestReleasedMovies), arg0, arg1)
}

// GetTheMostViewedMovies mocks base method.
func (m *MockStore) GetTheMostViewedMovies(arg0 context.Context, arg1 mongo0.GetMoviesParams) ([]mongo0.Movies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTheMostViewedMovies", arg0, arg1)
	ret0, _ := ret[0].([]mongo0.Movies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTheMostViewedMovies indicates an expected call of GetTheMostViewedMovies.
func (mr *MockStoreMockRecorder) GetTheMostViewedMovies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTheMostViewedMovies", reflect.TypeOf((*MockStore)(nil).GetTheMostViewedMovies), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockStore) GetUserByEmail(arg0 context.Context, arg1 string) (mongo0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(mongo0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockStoreMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockStore)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockStore) GetUserByID(arg0 context.Context, arg1 primitive.ObjectID) (mongo0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(mongo0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStoreMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStore)(nil).GetUserByID), arg0, arg1)
}

// GetUserByName mocks base method.
func (m *MockStore) GetUserByName(arg0 context.Context, arg1 string) (mongo0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", arg0, arg1)
	ret0, _ := ret[0].(mongo0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockStoreMockRecorder) GetUserByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockStore)(nil).GetUserByName), arg0, arg1)
}

// ReplaceMovieInfoByID mocks base method.
func (m *MockStore) ReplaceMovieInfoByID(arg0 context.Context, arg1 primitive.ObjectID, arg2 mongo0.Movies) (*mongo.UpdateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceMovieInfoByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*mongo.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceMovieInfoByID indicates an expected call of ReplaceMovieInfoByID.
func (mr *MockStoreMockRecorder) ReplaceMovieInfoByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceMovieInfoByID", reflect.TypeOf((*MockStore)(nil).ReplaceMovieInfoByID), arg0, arg1, arg2)
}

// SearchForMovies mocks base method.
func (m *MockStore) SearchForMovies(arg0 context.Context, arg1 mongo0.SearchForMoviesParams) ([]mongo0.Movies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchForMovies", arg0, arg1)
	ret0, _ := ret[0].([]mongo0.Movies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchForMovies indicates an expected call of SearchForMovies.
func (mr *MockStoreMockRecorder) SearchForMovies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchForMovies", reflect.TypeOf((*MockStore)(nil).SearchForMovies), arg0, arg1)
}

// UpdateComment mocks base method.
func (m *MockStore) UpdateComment(arg0 context.Context, arg1 mongo0.Comments) (*mongo.UpdateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComment", arg0, arg1)
	ret0, _ := ret[0].(*mongo.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateComment indicates an expected call of UpdateComment.
func (mr *MockStoreMockRecorder) UpdateComment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockStore)(nil).UpdateComment), arg0, arg1)
}

// UpdateUserName mocks base method.
func (m *MockStore) UpdateUserName(arg0 context.Context, arg1 mongo0.User) (*mongo.UpdateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserName", arg0, arg1)
	ret0, _ := ret[0].(*mongo.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserName indicates an expected call of UpdateUserName.
func (mr *MockStoreMockRecorder) UpdateUserName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserName", reflect.TypeOf((*MockStore)(nil).UpdateUserName), arg0, arg1)
}

// UpdateUserPassword mocks base method.
func (m *MockStore) UpdateUserPassword(arg0 context.Context, arg1 mongo0.User) (*mongo.UpdateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", arg0, arg1)
	ret0, _ := ret[0].(*mongo.UpdateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockStoreMockRecorder) UpdateUserPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockStore)(nil).UpdateUserPassword), arg0, arg1)
}
