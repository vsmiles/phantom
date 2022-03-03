package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	db "phantom/db/mongo"
)

var errMovieNotFound error = errors.New("movie is not found")

type createMovieRequest struct {
	Plot             string             `json:"plot"`
	Genres           []string           `json:"genres"`
	Runtime          int64              `json:"runtime"`
	Rated            string             `json:"rated"`
	Cast             []string           `json:"cast"`
	NumMflixComments int64              `json:"num_mflix_comments"`
	Poster           string             `json:"poster"`
	Title            string             `json:"title" binding:"required,min=1"`
	Fullplot         string             `json:"fullplot"`
	Languages        []string           `json:"languages"`
	Released         primitive.DateTime `json:"released"`
	Directors        []string           `json:"directors"`
	Writers          []string           `json:"writers"`
	Awards           struct {
		Wins        int64  `json:"wins"`
		Nominations int64  `json:"nominations"`
		Text        string `json:"text"`
	} `json:"awards"`
	Lastupdated string `json:"lastupdated"`
	Year        int64  `json:"year"`
	Imdb        struct {
		Rating float64 `json:"rating"`
		Votes  int64   `json:"votes"`
		Id     int64   `json:"id"`
	} `json:"imdb"`
	Countries []string `json:"countries"`
	Type      string   `json:"type"`
	Tomatoes  struct {
		Viewer struct {
			Rating     float64 `json:"rating"`
			NumReviews int64   `json:"numReviews"`
			Meter      int64   `json:"meter"`
		} `json:"viewer"`
		LastUpdated primitive.DateTime `json:"lastUpdated"`
	} `json:"tomatoes"`
}

// TODO: addOneMovieInfo the function is not yet complete
func (server *Server) createMovie(ctx *gin.Context) {
	var req createMovieRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.AddMovieParams{
		Genres:   req.Genres,
		Runtime:  req.Runtime,
		Title:    req.Title,
		Released: req.Released,
		Year:     req.Year,
	}
	id, err := server.store.AddMovie(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id.Hex()})
}

type getMoviesResponse struct {
	Id     string   `json:"id"`
	Plot   string   `json:"plot"`
	Genres []string `json:"genres"`
	Cast   []string `json:"cast"`
	Poster string   `json:"poster"`
	Title  string   `json:"title"`
	Year   int64    `json:"year"`
	Imdb   struct {
		Rating float64 `json:"rating"`
		Votes  int64   `json:"votes"`
		Id     int64   `json:"id"`
	} `json:"imdb"`
	Countries []string `json:"countries"`
}

type movieIdRequest struct {
	Id string `uri:"id" binding:"required,hexadecimal,min=24"`
}

// getMovie  is to get the details of a movie by its id
func (server *Server) getMovie(ctx *gin.Context) {
	var req movieIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	movie, err := server.store.GetMovieByID(ctx, objectID)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			ctx.JSON(http.StatusNotFound, errorResponse(errMovieNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movie)
}

type searchForMoviesRequest struct {
	Search   string `form:"search" binding:"required,min=1"`
	PageSize int64  `form:"s" binding:"required,max=50"`
	PageId   int64  `form:"p" binding:"required,min=1"`
}

// searchForMovies is to retrieve movie information
func (server *Server) searchForMovies(ctx *gin.Context) {
	var req searchForMoviesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.SearchForMoviesParams{
		Text:  req.Search,
		Skip:  req.PageSize * (req.PageId - 1),
		Limit: req.PageSize,
	}
	movies, err := server.store.SearchForMovies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := generateGetMoviesResponse(movies)
	ctx.JSON(http.StatusOK, rsp)
}

type listMoviesByGenresRequest struct {
	Genres   string `form:"genres" binding:"required,min=1"`
	Sort     string `form:"sort" binding:"required,oneof=hotness time rating"`
	PageSize int64  `form:"s" binding:"required,min=1,max=50"`
	PageId   int64  `form:"p" binding:"required,min=1"`
}

// listMoviesByGenres is to get some movie information according to genre
func (server *Server) listMoviesByGenres(ctx *gin.Context) {
	var req listMoviesByGenresRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.GetMoviesParams{
		Genres:      req.Genres,
		SortOptions: "sort_" + req.Sort,
		Skip:        req.PageSize * (req.PageId - 1),
		Limit:       req.PageSize,
	}
	movies, err := server.store.GetMoviesByGenres(ctx, arg)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("movie is not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := generateGetMoviesResponse(movies)
	ctx.JSON(http.StatusOK, rsp)
}

type getMoviesRequest struct {
	PageSize int64 `form:"s" binding:"required,min=1,max=50"`
	PageId   int64 `form:"p" binding:"required,min=1"`
}

func (server *Server) listTheMostWatchedMovies(ctx *gin.Context) {
	var req getMoviesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.GetMoviesParams{
		Skip:  req.PageSize * (req.PageId - 1),
		Limit: req.PageSize,
	}
	movies, err := server.store.GetTheMostViewedMovies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := generateGetMoviesResponse(movies)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listTheLatestReleasedMovies(ctx *gin.Context) {
	var req getMoviesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.GetMoviesParams{
		Skip:  req.PageSize * (req.PageId - 1),
		Limit: req.PageSize,
	}
	movies, err := server.store.GetTheLatestReleasedMovies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := generateGetMoviesResponse(movies)
	ctx.JSON(http.StatusOK, rsp)
}

// getMoviesWithMostCommented TODO: This function is used to list the most commented movies from users
func (server *Server) getTheMostCommentedMovies() {

}

// TODO: 预览播放功能
// recommendedMovies TODO : This function is used to recommend movies for users
func (server *Server) recommendedMovies(ctx *gin.Context) {

}

func (server *Server) updateMovie(ctx *gin.Context) {
	var reqUri movieIdRequest
	var reqJson createMovieRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(reqUri.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.Movies{
		Genres:   reqJson.Genres,
		Runtime:  reqJson.Runtime,
		Title:    reqJson.Title,
		Released: reqJson.Released,
		Year:     reqJson.Year,
	}
	_, err = server.store.ReplaceMovieInfoByID(ctx, objectId, arg)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			ctx.JSON(http.StatusNotFound, errorResponse(errMovieNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"updated": "OK"})
}

func generateGetMoviesResponse(movies []db.Movies) []getMoviesResponse {
	var rsp []getMoviesResponse
	for _, movie := range movies {
		arg := getMoviesResponse{
			Id:     movie.Id.Hex(),
			Plot:   movie.Plot,
			Genres: movie.Genres,
			Cast:   movie.Cast,
			Poster: movie.Poster,
			Title:  movie.Title,
			Year:   movie.Year,
			Imdb: struct {
				Rating float64 `json:"rating"`
				Votes  int64   `json:"votes"`
				Id     int64   `json:"id"`
			}(movie.Imdb),
			Countries: movie.Countries,
		}
		rsp = append(rsp, arg)
	}
	return rsp
}
