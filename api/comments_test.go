package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	mockdb "phantom/db/mock"
	db "phantom/db/mongo"
	"phantom/token"
	"phantom/util"
	"testing"
	"time"
)

func TestCreateCommentAPI(t *testing.T) {
	returnId := primitive.NewObjectID()
	comment := randomComment()
	testCase := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    comment.Email,
				"movie_id": comment.MovieID.Hex(),
				"text":     comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddCommentParams{
					Name:    "user",
					Email:   comment.Email,
					MovieID: comment.MovieID,
					Text:    comment.Text,
				}
				store.EXPECT().AddComment(gomock.Any(), gomock.Eq(arg)).Times(1).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchObjectId(t, returnId, recorder.Body)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"email":    "Invalid",
				"movie_id": comment.MovieID,
				"text":     comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddComment(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidMovieID_Length",
			body: gin.H{
				"email":    comment.Email,
				"movie_id": comment.MovieID[2:],
				"text":     comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddComment(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidText",
			body: gin.H{
				"email":    comment.Email,
				"movie_id": comment.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddComment(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidMovieId_ObjectId",
			body: gin.H{
				"email":    comment.Email,
				"movie_id": util.RandomString(24),
				"text":     comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddComment(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnAuthorized",
			body: gin.H{
				"email":    comment.Email,
				"movie_id": comment.MovieID.Hex(),
				"text":     comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddComment(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"email":    comment.Email,
				"movie_id": comment.MovieID.Hex(),
				"text":     comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddCommentParams{
					Name:    "user",
					Email:   comment.Email,
					MovieID: comment.MovieID,
					Text:    comment.Text,
				}
				store.EXPECT().AddComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnId, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/comments"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListCommentsByMovieIDAPI(t *testing.T) {
	n := 5
	var returnComments []db.Comments
	movieId := primitive.NewObjectID()
	for i := 0; i < n; i++ {
		comment := db.Comments{
			Name:    util.RandomUser(),
			MovieID: movieId,
			Text:    util.RandomString(10),
		}
		returnComments = append(returnComments, comment)
	}
	testCase := []struct {
		name          string
		query         listCommentsRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listCommentsRequest{
				MovieID:  movieId.Hex(),
				PageSize: int64(n),
				PageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetCommentsParams{
					MovieID: movieId,
					Skip:    5,
					Limit:   5,
				}
				store.EXPECT().GetCommentsByMovieID(gomock.Any(), gomock.Eq(arg)).Times(1).Return(returnComments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchComments(t, returnComments, recorder.Body)
			},
		},
		{
			name: "InvalidMovieId_Length",
			query: listCommentsRequest{
				MovieID:  movieId.Hex()[2:],
				PageSize: int64(n),
				PageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommentsByMovieID(gomock.Any(), gomock.Any()).Times(0).Return(returnComments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidMovieId_ObjectId",
			query: listCommentsRequest{
				MovieID:  util.RandomString(24),
				PageSize: int64(n),
				PageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommentsByMovieID(gomock.Any(), gomock.Any()).Times(0).Return(returnComments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: listCommentsRequest{
				MovieID:  movieId.Hex(),
				PageSize: int64(100),
				PageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommentsByMovieID(gomock.Any(), gomock.Any()).Times(0).Return(returnComments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageId",
			query: listCommentsRequest{
				MovieID:  movieId.Hex(),
				PageSize: int64(n),
				PageId:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommentsByMovieID(gomock.Any(), gomock.Any()).Times(0).Return(returnComments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: listCommentsRequest{
				MovieID:  movieId.Hex(),
				PageSize: int64(n),
				PageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetCommentsParams{
					MovieID: movieId,
					Skip:    5,
					Limit:   5,
				}
				store.EXPECT().GetCommentsByMovieID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnComments, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/comments"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("movie_id", tc.query.MovieID)
			q.Add("s", fmt.Sprintf("%v", tc.query.PageSize))
			q.Add("p", fmt.Sprintf("%v", tc.query.PageId))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCommentAPI(t *testing.T) {
	comment := db.Comments{
		ID:   primitive.NewObjectID(),
		Name: "user",
		Text: util.RandomString(10),
	}
	testCase := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"id":   comment.ID.Hex(),
				"text": comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.Comments{
					ID:   comment.ID,
					Name: "user",
					Text: comment.Text,
				}
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidId_Length",
			body: gin.H{
				"id":   comment.ID.Hex()[2:],
				"text": comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidId_ObjectId",
			body: gin.H{
				"id":   util.RandomString(24),
				"text": comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidText",
			body: gin.H{
				"id": comment.ID.Hex(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnAuthorized",
			body: gin.H{
				"id":   comment.ID.Hex(),
				"text": comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id":   comment.ID.Hex(),
				"text": comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.Comments{
					ID:   comment.ID,
					Name: "user",
					Text: comment.Text,
				}
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"id":   comment.ID.Hex(),
				"text": comment.Text,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.Comments{
					ID:   comment.ID,
					Name: "user",
					Text: comment.Text,
				}
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/comments"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCommentAPI(t *testing.T) {
	commentId := primitive.NewObjectID()
	testCase := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"id": commentId.Hex(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Eq(commentId), gomock.Eq("user")).
					Times(1).
					Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidId_Length",
			body: gin.H{
				"id": commentId.Hex()[2:],
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidId_ObjectId",
			body: gin.H{
				"id": util.RandomString(24),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnAuthorized",
			body: gin.H{
				"id": commentId.Hex(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"id": commentId.Hex(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Eq(commentId), gomock.Eq("user")).
					Times(1).
					Return(int64(0), mongo.ErrNoDocuments)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id": commentId.Hex(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Eq(commentId), gomock.Eq("user")).
					Times(1).
					Return(int64(0), mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/comments"
			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomComment() db.AddCommentParams {
	return db.AddCommentParams{
		Name:    util.RandomUser(),
		Email:   util.RandomEmail(),
		MovieID: primitive.NewObjectID(),
		Text:    util.RandomString(10),
	}

}

func requireBodyMatchComments(t *testing.T, returnComments []db.Comments, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	var comments []ListCommentsResponse
	json.Unmarshal(data, &comments)
	for i := range comments {
		require.Equal(t, returnComments[i].Name, comments[i].Name)
	}

}
