package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Movies struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Plot             string             `json:"plot" bson:"plot,omitempty"`
	Genres           []string           `json:"genres" bson:"genres,omitempty"`
	Runtime          int64              `json:"runtime" bson:"runtime,omitempty"`
	Rated            string             `json:"rated" bson:"rated,omitempty"`
	Cast             []string           `json:"cast" bson:"cast,omitempty"`
	NumMflixComments int64              `json:"num_mflix_comments" bson:"num_mflix_comments,omitempty"`
	Poster           string             `json:"poster" bson:"poster,omitempty"`
	Title            string             `json:"title" bson:"title,omitempty"`
	Fullplot         string             `json:"fullplot" bson:"fullplot,omitempty"`
	Languages        []string           `json:"languages" bson:"languages,omitempty"`
	Released         primitive.DateTime `json:"released" bson:"released,omitempty"`
	Directors        []string           `json:"directors" bson:"directors,omitempty"`
	Writers          []string           `json:"writers" bson:"writers,omitempty"`
	Awards           struct {
		Wins        int64  `json:"wins" bson:"wins,omitempty"`
		Nominations int64  `json:"nominations" bson:"nominations,omitempty"`
		Text        string `json:"text" bson:"text,omitempty"`
	} `json:"awards" bson:"awards,omitempty"`
	Lastupdated string `json:"lastupdated" bson:"lastupdated,omitempty"`
	Year        int64  `json:"year" bson:"year,omitempty"`
	Imdb        struct {
		Rating float64 `json:"rating" bson:"rating,omitempty"`
		Votes  int64   `json:"votes" bson:"votes,omitempty"`
		Id     int64   `json:"id" bson:"id,omitempty"`
	} `json:"imdb" bson:"imdb,omitempty"`
	Countries []string `json:"countries" bson:"countries,omitempty"`
	Type      string   `json:"type" bson:"type,omitempty"`
	Tomatoes  struct {
		Viewer struct {
			Rating     float64 `json:"rating" bson:"rating,omitempty"`
			NumReviews int64   `json:"numReviews" bson:"numReviews,omitempty"`
			Meter      int64   `json:"meter" bson:"meter,omitempty"`
		} `json:"viewer" bson:"viewer,omitempty"`
		LastUpdated primitive.DateTime `json:"lastUpdated" bson:"lastUpdated,omitempty"`
	} `json:"tomatoes" bson:"tomatoes,omitempty"`
}

type AddMovieParams struct {
	Genres   []string           `json:"genres" bson:"genres,omitempty"`
	Runtime  int64              `json:"runtime" bson:"runtime,omitempty"`
	Title    string             `json:"title" bson:"title,omitempty"`
	Released primitive.DateTime `json:"released" bson:"released,omitempty"`
	Year     int64              `json:"year" bson:"year,omitempty"`
	Imdb     struct {
		Rating float64 `json:"rating" bson:"rating,omitempty"`
	} `json:"imdb" bson:"imdb,omitempty"`
}

// AddMovie can add a movie information
func (q *Queries) AddMovie(ctx context.Context, arg AddMovieParams) (primitive.ObjectID, error) {
	res, err := q.movies.InsertOne(ctx, arg)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

// AddMovies can add a set of movie information
func (q *Queries) AddMovies(ctx context.Context, arg []AddMovieParams) ([]interface{}, error) {
	movies := make([]interface{}, len(arg))
	for i, v := range arg {
		movies[i] = v
	}
	res, err := q.movies.InsertMany(ctx, movies)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

// GetMovieByID can get the movie information by movie id
func (q *Queries) GetMovieByID(ctx context.Context, id primitive.ObjectID) (Movies, error) {
	var movie Movies
	err := q.movies.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		return Movies{}, err
	}
	return movie, nil
}

type GetMoviesParams struct {
	Title       string `json:"title" bson:"title,omitempty"`
	Genres      string `json:"genres"`
	SortOptions string `json:"sort_options"`
	Skip        int64  `json:"skip"`
	Limit       int64  `json:"limit"`
}

func (q *Queries) SearchForMovies(ctx context.Context, arg SearchForMoviesParams) ([]Movies, error) {
	matchStage := bson.D{{"$match", bson.D{{"$text", bson.D{{"$search", arg.Text}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"score", bson.D{{"$meta", "textScore"}}}}}}
	projectStage := projectStage()
	skipStage := bson.D{{"$skip", arg.Skip}}
	limitStage := bson.D{{"$limit", arg.Limit}}
	cursor, err := q.movies.Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, projectStage, skipStage, limitStage})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var movies []Movies
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

// GetMoviesByGenres Get the movie information of the past year by movie genres,
// you can also choose to sort by hotness, sort by time, sort by rating,
// the default is sort by hotness
func (q *Queries) GetMoviesByGenres(ctx context.Context, arg GetMoviesParams) ([]Movies, error) {
	projectStage := projectStage()
	matchStage := bson.D{{
		"$match", bson.D{
			{"genres", arg.Genres},
			{"released", bson.D{
				{"$gte", primitive.NewDateTimeFromTime(time.Now().AddDate(-1, 0, 0))},
			}}}}}
	sortStageWithRuntime := bson.D{{"$sort", bson.D{{"runtime", -1}}}}
	sortStageWithTime := bson.D{{"$sort", bson.D{{"released", -1}}}}
	sortStageWithRating := bson.D{{"$sort", bson.D{{"imdb.rating", -1}}}}
	skipStage := bson.D{{"$skip", arg.Skip}}
	limitStage := bson.D{{"$limit", arg.Limit}}
	sortStage := sortStageWithRuntime
	switch arg.SortOptions {
	case "sort_time":
		sortStage = sortStageWithTime
	case "sort_rating":
		sortStage = sortStageWithRating
	}
	cursor, err := q.movies.Aggregate(
		ctx,
		mongo.Pipeline{
			projectStage,
			matchStage,
			sortStage,
			skipStage,
			limitStage,
		},
	)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var movies []Movies
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

type SearchForMoviesParams struct {
	Text  string `json:"text"`
	Skip  int64  `json:"skip"`
	Limit int64  `json:"limit"`
}

func (q *Queries) GetTheMostViewedMovies(ctx context.Context, arg GetMoviesParams) ([]Movies, error) {
	projectStage := projectStage()
	sortStage := bson.D{{"$sort", bson.D{{"runtime", -1}}}}
	skipStage := bson.D{{"$skip", arg.Skip}}
	limitStage := bson.D{{"$limit", arg.Limit}}
	cursor, err := q.movies.Aggregate(ctx, mongo.Pipeline{sortStage, projectStage, skipStage, limitStage})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var movies []Movies
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

// GetTheLatestReleasedMovies is to get the latest released movies
func (q *Queries) GetTheLatestReleasedMovies(ctx context.Context, arg GetMoviesParams) ([]Movies, error) {
	projectStage := projectStage()
	sortStage := bson.D{{"$sort", bson.D{{"released", -1}}}}
	skipStage := bson.D{{"$skip", arg.Skip}}
	limitStage := bson.D{{"$limit", arg.Limit}}
	cursor, err := q.movies.Aggregate(ctx, mongo.Pipeline{projectStage, sortStage, skipStage, limitStage})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var movies []Movies
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

// GetTheMostCommentedMovies TODO: This function is used to get the most commented movies from users
func (q *Queries) GetTheMostCommentedMovies(ctx context.Context, arg GetMoviesParams) ([]Movies, error) {
	return nil, nil
}

// GetTheMostRecentlyMovies TODO: This function is used to get the most recently uploaded movies
func (q *Queries) GetTheMostRecentlyMovies(ctx context.Context, arg GetMoviesParams) ([]Movies, error) {
	return nil, nil
}

func (q *Queries) ReplaceMovieInfoByID(ctx context.Context, id primitive.ObjectID, movie Movies) (*mongo.UpdateResult, error) {
	data, err := bson.Marshal(&movie)
	if err != nil {
		return nil, err
	}
	updateResult, err := q.movies.ReplaceOne(ctx, bson.M{"_id": id}, data)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return updateResult, nil
}

func (q *Queries) DeleteMovieByID(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := q.movies.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	if deleteResult.DeletedCount == 0 {
		return 0, mongo.ErrNoDocuments
	}
	return deleteResult.DeletedCount, nil
}

func projectStage() bson.D {
	return bson.D{{
		"$project", bson.D{
			{"_id", 1},
			{"plot", 1},
			{"genres", 1},
			{"runtime", 1},
			{"cast", 1},
			{"poster", 1},
			{"title", 1},
			{"released", 1},
			{"year", 1},
			{"imdb", 1},
			{"countries", 1},
		},
	}}
}
