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
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestCreateMovieApi(t *testing.T) {
	movie := randomMovie()
	returnId := primitive.NewObjectID()
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
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddMovie(gomock.Any(), gomock.Eq(movie)).Times(1).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchObjectId(t, returnId, recorder.Body)
			},
		},
		{
			name: "InvalidTitle",
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddMovie(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddMovie(gomock.Any(), gomock.Eq(movie)).
					Times(1).
					Return(returnId, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UnAuthorized",
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddMovie(gomock.Any(), gomock.Any()).Times(0).Return(returnId, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			url := "/movies"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestGetMovieApi(t *testing.T) {
	objectId := primitive.NewObjectID()
	returnMovie := db.Movies{
		Id: objectId,
	}
	testCase := []struct {
		name          string
		movieId       string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieId: objectId.Hex(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovieByID(gomock.Any(), gomock.Eq(objectId)).Times(1).Return(returnMovie, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMovieDetails(t, returnMovie, recorder.Body)
			},
		},
		{
			name:    "InvalidId_Length",
			movieId: objectId.Hex()[2:],
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovieByID(gomock.Any(), gomock.Any()).Times(0).Return(returnMovie, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "InvalidId_Type",
			movieId: util.RandomString(24),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovieByID(gomock.Any(), gomock.Any()).Times(0).Return(returnMovie, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			movieId: objectId.Hex(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovieByID(gomock.Any(), gomock.Eq(objectId)).
					Times(1).
					Return(returnMovie, mongo.ErrNoDocuments)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			movieId: objectId.Hex(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovieByID(gomock.Any(), gomock.Eq(objectId)).
					Times(1).
					Return(returnMovie, mongo.ErrClientDisconnected)
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
			url := fmt.Sprintf("/movies/%s", tc.movieId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestSearchForMoviesApi(t *testing.T) {
	n := 5
	var returnMovies []db.Movies
	genres := []string{util.RandomString(6), util.RandomString(6)}
	for i := 0; i < n; i++ {
		movie := db.Movies{
			Id:     primitive.NewObjectID(),
			Genres: genres,
		}
		returnMovies = append(returnMovies, movie)
	}

	type Query struct {
		Search   string
		pageSize int64
		pageId   int64
	}
	testCase := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				Search:   genres[0],
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.SearchForMoviesParams{
					Text:  genres[0],
					Skip:  5,
					Limit: 5,
				}
				store.EXPECT().SearchForMovies(gomock.Any(), gomock.Eq(arg)).Times(1).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchManyMovies(t, returnMovies, recorder.Body)
			},
		},
		{
			name: "InvalidText",
			query: Query{
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().SearchForMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				Search:   genres[0],
				pageSize: int64(100),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().SearchForMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageId",
			query: Query{
				Search:   genres[0],
				pageSize: int64(n),
				pageId:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().SearchForMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				Search:   genres[0],
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.SearchForMoviesParams{
					Text:  genres[0],
					Skip:  5,
					Limit: 5,
				}
				store.EXPECT().SearchForMovies(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnMovies, mongo.ErrClientDisconnected)
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
			url := "/search"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request utl
			q := request.URL.Query()
			q.Add("search", tc.query.Search)
			q.Add("s", fmt.Sprintf("%v", tc.query.pageSize))
			q.Add("p", fmt.Sprintf("%v", tc.query.pageId))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}

func TestListMoviesByGenresAPI(t *testing.T) {
	n := 5
	var returnMovies []db.Movies
	genres := []string{util.RandomString(6), util.RandomString(6)}
	for i := 0; i < n; i++ {
		movie := db.Movies{
			Id:     primitive.NewObjectID(),
			Genres: genres,
		}
		returnMovies = append(returnMovies, movie)
	}
	sort.Slice(returnMovies, func(i, j int) bool {
		return returnMovies[i].Imdb.Rating > returnMovies[j].Imdb.Rating // Desc
	})
	type Query struct {
		genres   string
		sort     string
		pageSize int64
		pageId   int64
	}
	testCase := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				genres:   genres[0],
				sort:     "rating",
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Genres:      genres[0],
					SortOptions: "sort_rating",
					Skip:        5,
					Limit:       5,
				}
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Eq(arg)).Times(1).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchManyMovies(t, returnMovies, recorder.Body)
			},
		},
		{
			name: "InvalidGenres",
			query: Query{
				sort:     "rating",
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidSort",
			query: Query{
				genres:   genres[0],
				sort:     "Invalid",
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				genres:   genres[0],
				sort:     "rating",
				pageSize: int64(100),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageId",
			query: Query{
				genres:   genres[0],
				sort:     "rating",
				pageSize: int64(n),
				pageId:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			query: Query{
				genres:   genres[0],
				sort:     "rating",
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Genres:      genres[0],
					SortOptions: "sort_rating",
					Skip:        5,
					Limit:       5,
				}
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnMovies, mongo.ErrNoDocuments)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				genres:   genres[0],
				sort:     "rating",
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Genres:      genres[0],
					SortOptions: "sort_rating",
					Skip:        5,
					Limit:       5,
				}
				store.EXPECT().GetMoviesByGenres(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnMovies, mongo.ErrClientDisconnected)
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
			url := "/movies/genres"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request utl
			q := request.URL.Query()
			q.Add("genres", tc.query.genres)
			q.Add("sort", tc.query.sort)
			q.Add("s", fmt.Sprintf("%v", tc.query.pageSize))
			q.Add("p", fmt.Sprintf("%v", tc.query.pageId))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}

func TestListTheMostWatchedMoviesApi(t *testing.T) {
	n := 5
	var returnMovies []db.Movies
	for i := 0; i < n; i++ {
		movie := db.Movies{
			Id:      primitive.NewObjectID(),
			Runtime: util.RandomInt(1, 100),
		}
		returnMovies = append(returnMovies, movie)
	}
	sort.Slice(returnMovies, func(i, j int) bool {
		return returnMovies[i].Runtime > returnMovies[j].Runtime // Desc
	})
	type Query struct {
		pageSize int64
		pageId   int64
	}
	testCase := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Skip:  5,
					Limit: 5,
				}
				store.EXPECT().GetTheMostViewedMovies(gomock.Any(), gomock.Eq(arg)).Times(1).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchManyMovies(t, returnMovies, recorder.Body)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageSize: int64(100),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTheMostViewedMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageId",
			query: Query{
				pageSize: int64(n),
				pageId:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTheMostViewedMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Skip:  5,
					Limit: 5,
				}
				store.EXPECT().GetTheMostViewedMovies(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnMovies, mongo.ErrClientDisconnected)
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
			url := "/movies/most_watched"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request utl
			q := request.URL.Query()
			q.Add("s", strconv.FormatInt(tc.query.pageSize, 10))
			q.Add("p", strconv.FormatInt(tc.query.pageId, 10))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}

func TestListTheLatestReleasedApi(t *testing.T) {
	n := 5
	var returnMovies []db.Movies
	for i := 0; i < n; i++ {
		movie := db.Movies{
			Id:       primitive.NewObjectID(),
			Released: primitive.NewDateTimeFromTime(time.Now()),
		}
		returnMovies = append(returnMovies, movie)
	}
	sort.Slice(returnMovies, func(i, j int) bool {
		return returnMovies[i].Released > returnMovies[j].Released
	})

	type Query struct {
		pageSize int64
		pageId   int64
	}
	testCase := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Skip:  5,
					Limit: 5,
				}
				store.EXPECT().GetTheLatestReleasedMovies(gomock.Any(), gomock.Eq(arg)).Times(1).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchManyMovies(t, returnMovies, recorder.Body)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageSize: int64(100),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTheLatestReleasedMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageId",
			query: Query{
				pageSize: int64(n),
				pageId:   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTheLatestReleasedMovies(gomock.Any(), gomock.Any()).Times(0).Return(returnMovies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageSize: int64(n),
				pageId:   2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetMoviesParams{
					Skip:  5,
					Limit: 5,
				}
				store.EXPECT().GetTheLatestReleasedMovies(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnMovies, mongo.ErrClientDisconnected)
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
			url := "/movies/Latest"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request utl
			q := request.URL.Query()
			q.Add("s", fmt.Sprintf("%v", tc.query.pageSize))
			q.Add("p", fmt.Sprintf("%v", tc.query.pageId))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}

func TestUpdateMovieApi(t *testing.T) {
	movie := randomMovie()
	ObjectId := primitive.NewObjectID()
	testCase := []struct {
		name          string
		id            string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   ObjectId.Hex(),
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.Movies{
					Genres:   movie.Genres,
					Runtime:  movie.Runtime,
					Title:    movie.Title,
					Released: movie.Released,
					Year:     movie.Year,
				}
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Eq(ObjectId), gomock.Eq(arg)).
					Times(1).
					Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidId_Length",
			id:   ObjectId.Hex()[2:],
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidId_type",
			id:   util.RandomString(24),
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			id:   ObjectId.Hex(),
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   ObjectId.Hex(),
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.Movies{
					Genres:   movie.Genres,
					Runtime:  movie.Runtime,
					Title:    movie.Title,
					Released: movie.Released,
					Year:     movie.Year,
				}
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Eq(ObjectId), gomock.Eq(arg)).
					Times(1).
					Return(nil, mongo.ErrClientDisconnected)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			id:   ObjectId.Hex(),
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.Movies{
					Genres:   movie.Genres,
					Runtime:  movie.Runtime,
					Title:    movie.Title,
					Released: movie.Released,
					Year:     movie.Year,
				}
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Eq(ObjectId), gomock.Eq(arg)).
					Times(1).
					Return(nil, mongo.ErrNoDocuments)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "UnAuthorized",
			id:   ObjectId.Hex(),
			body: gin.H{
				"genres":   movie.Genres,
				"runtime":  movie.Runtime,
				"title":    movie.Title,
				"released": movie.Released,
				"year":     movie.Year,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ReplaceMovieInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			url := fmt.Sprintf("/movies/%s", tc.id)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchManyMovies(t *testing.T, movies []db.Movies, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotMovies []getMoviesResponse
	err = json.Unmarshal(data, &gotMovies)
	require.NoError(t, err)
	for i := range gotMovies {
		require.Equal(t, movies[i].Id.Hex(), gotMovies[i].Id)
	}

}

func requireBodyMatchObjectId(t *testing.T, objectId primitive.ObjectID, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	var gotMovie getMoviesResponse
	err = json.Unmarshal(data, &gotMovie)
	require.NoError(t, err)
	require.Equal(t, objectId.Hex(), gotMovie.Id)
}

func requireBodyMatchMovieDetails(t *testing.T, movie db.Movies, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	var gotMovie db.Movies
	err = json.Unmarshal(data, &gotMovie)
	require.NoError(t, err)
	require.Equal(t, movie.Id, gotMovie.Id)
}

func randomMovie() db.AddMovieParams {
	date := util.RandomDate()
	return db.AddMovieParams{
		Runtime:  util.RandomInt(0, 100),
		Title:    util.RandomString(8),
		Released: primitive.NewDateTimeFromTime(date),
		Genres:   util.RandomGenres(),
		Year:     int64(date.Year()),
	}
}
