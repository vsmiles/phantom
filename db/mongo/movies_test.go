package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"phantom/util"
	"sort"
	"testing"
	"time"
)

func randomMovie() AddMovieParams {
	date := util.RandomDate()
	return AddMovieParams{
		Runtime:  util.RandomInt(0, 100),
		Title:    util.RandomString(8),
		Released: primitive.NewDateTimeFromTime(date),
		Genres:   util.RandomGenres(),
		Year:     int64(date.Year()),
		Imdb: struct {
			Rating float64 `json:"rating" bson:"rating,omitempty"`
		}{Rating: util.Randomfloat(1, 10)},
	}
}

func addMovie(t *testing.T, movie AddMovieParams) primitive.ObjectID {
	id, err := testQueries.AddMovie(context.Background(), movie)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	return id
}

func TestAddMovie(t *testing.T) {
	addMovie(t, randomMovie())
}

func getMovieByID(t *testing.T, id primitive.ObjectID) Movies {
	movie, err := testQueries.GetMovieByID(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, movie)
	return movie
}

func TestGetMovieByID(t *testing.T) {
	movie1 := randomMovie()
	id := addMovie(t, movie1)
	movie2 := getMovieByID(t, id)
	require.Equal(t, movie1.Runtime, movie2.Runtime)
	require.Equal(t, movie1.Title, movie2.Title)
	require.Equal(t, movie1.Released, movie2.Released)
	require.Equal(t, movie1.Genres, movie2.Genres)
}

func TestSearchForMovies(t *testing.T) {
	movie := randomMovie()
	addMovie(t, movie)
	arg := SearchForMoviesParams{
		Text:  movie.Title,
		Skip:  0,
		Limit: 5,
	}
	movies, err := testQueries.SearchForMovies(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, movie)
	require.GreaterOrEqual(t, len(movies), 1)
	require.LessOrEqual(t, len(movies), 5)
	arg.Text = util.RandomString(9)
}

func TestGetMoviesByGenres(t *testing.T) {
	n := 10
	genres := []string{util.RandomString(8), util.RandomString(8)}
	var movies1 []AddMovieParams
	for i := 0; i < n; i++ {
		movie := randomMovie()
		movie.Genres = genres
		movies1 = append(movies1, movie)
		addMovie(t, movie)
	}
	log.Println("movie year:", movies1[0].Year)

	arg := GetMoviesParams{
		Genres:      genres[0],
		SortOptions: "sort_hotness",
		Skip:        0,
		Limit:       5,
	}
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			arg.SortOptions = "sort_hotness"
			sort.Slice(movies1, func(i, j int) bool {
				return movies1[i].Runtime > movies1[j].Runtime // desc
			})
		case 1:
			arg.SortOptions = "sort_time"
			sort.Slice(movies1, func(i, j int) bool {
				return movies1[i].Released > movies1[j].Released // desc
			})
		case 2:
			arg.SortOptions = "sort_rating"
			sort.Slice(movies1, func(i, j int) bool {
				return movies1[i].Imdb.Rating > movies1[j].Imdb.Rating // desc
			})
		}
		movies2, err := testQueries.GetMoviesByGenres(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, movies2)
		require.Equal(t, 5, len(movies2))
		for j := range movies2 {
			require.NotEmpty(t, movies2[j])
			require.Equal(t, genres[0], movies2[j].Genres[0])
			switch i {
			case 0:
				require.Equal(t, movies1[j].Runtime, movies2[j].Runtime)
			case 1:
				require.Equal(t, movies1[j].Released, movies2[j].Released)
			case 2:
				require.Equal(t, movies1[j].Imdb.Rating, movies2[j].Imdb.Rating)
			}
		}
	}
}

func TestGetTheMostViewedMovies(t *testing.T) {
	n := 10
	var runtime int64 = 1500
	for i := 0; i < n; i++ {
		movie := randomMovie()
		movie.Runtime = runtime
		addMovie(t, movie)
	}
	arg := GetMoviesParams{
		Skip:  5,
		Limit: 5,
	}
	movies, err := testQueries.GetTheMostViewedMovies(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, movies)
	require.Equal(t, 5, len(movies))
	for i := range movies {
		require.Equal(t, runtime, movies[i].Runtime)
	}
}

func TestGetTheLatestReleasedMovies(t *testing.T) {
	n := 5
	var movies1 []AddMovieParams
	for i := 0; i < n; i++ {
		movie := randomMovie()
		movie.Released = primitive.NewDateTimeFromTime(time.Now())
		movies1 = append(movies1, movie)
		addMovie(t, movie)
	}
	sort.Slice(movies1, func(i, j int) bool {
		return movies1[i].Released > movies1[j].Released // desc
	})
	arg := GetMoviesParams{
		Skip:  0,
		Limit: 5,
	}
	movies2, err := testQueries.GetTheLatestReleasedMovies(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, movies2)
	for i := range movies2 {
		require.Equal(t, movies1[i].Title, movies2[i].Title)
	}
}

func TestReplaceMovieInfoByID(t *testing.T) {
	movie1 := getMovieByID(t, addMovie(t, randomMovie()))
	date := util.RandomDate()
	movie1.Released = primitive.NewDateTimeFromTime(date)
	movie1.Year = int64(date.Year())
	movie1.Runtime = util.RandomInt(1, 100)
	updateResult, err := testQueries.ReplaceMovieInfoByID(context.Background(), movie1.Id, movie1)
	require.NoError(t, err)
	require.NotEmpty(t, updateResult)
	movie2 := getMovieByID(t, movie1.Id)
	require.Equal(t, movie1, movie2)

	movie2.Id = primitive.NewObjectID()
	_, err = testQueries.ReplaceMovieInfoByID(context.Background(), movie2.Id, movie1)
	require.Equal(t, mongo.ErrNoDocuments, err)
}

func deleteMovieByID(t *testing.T, id primitive.ObjectID) int64 {
	deleteCount, err := testQueries.DeleteMovieByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, int64(1), deleteCount)
	return deleteCount
}
func TestDeleteMovieByID(t *testing.T) {
	id := addMovie(t, randomMovie())
	deleteMovieByID(t, id)

	_, err := testQueries.DeleteMovieByID(context.Background(), primitive.NewObjectID())
	require.Equal(t, mongo.ErrNoDocuments, err)
}
