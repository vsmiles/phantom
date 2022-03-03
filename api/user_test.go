package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	mockdb "phantom/db/mock"
	db "phantom/db/mongo"
	"phantom/util"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.AddUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.AddUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(arg.Password, e.password)
	if err != nil {
		return false
	}
	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func EqCreateUserParams(arg db.AddUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// TODO: Duplicate error has not been tested
func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserParams{
					Name:  user.Name,
					Email: user.Email,
				}
				store.EXPECT().AddUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"name":     "invalid#",
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserParams{
					Name:  user.Name,
					Email: user.Email,
				}
				store.EXPECT().AddUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(0).
					Return(user.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"name":     user.Name,
				"email":    "email",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserParams{
					Name:  user.Name,
					Email: user.Email,
				}
				store.EXPECT().AddUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(0).
					Return(user.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": "short",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddUserParams{
					Name:  user.Name,
					Email: user.Email,
				}
				store.EXPECT().AddUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(0).
					Return(user.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": util.RandomString(6),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user.ID, mongo.ErrClientDisconnected)
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

			url := "/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid#",
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(0).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "InvalidEmail",
			body: gin.H{
				"username": user.Name,
				"email":    "email",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(0).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "InvalidPassword",
			body: gin.H{
				"username": user.Name,
				"email":    user.Email,
				"password": "short",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(0).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Name,
				"email":    user.Email,
				"password": "Incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NonUsernameAndEmail",
			body: gin.H{
				"password": user.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).Times(0).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UserNotFound_ByEmail",
			body: gin.H{
				"username": user.Name,
				"email":    user.Email,
				"password": user.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, mongo.ErrNoDocuments)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "UserNotFound_ByUsername",
			body: gin.H{
				"username": user.Name,
				"password": user.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByName(gomock.Any(), gomock.Eq(user.Name)).
					Times(1).
					Return(user, mongo.ErrNoDocuments)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := "/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Name:     util.RandomUser(),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
	}
	return
}
