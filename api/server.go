package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "phantom/db/mongo"
	"phantom/token"
	"phantom/util"
)

// Server servers HTTP request for our
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Registration binding tag
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("login", validLogin)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/register", server.register)
	router.POST("/login", server.login)

	router.GET("/movies/:id", server.getMovie)
	router.GET("/search", server.searchForMovies)
	router.GET("/movies/genres", server.listMoviesByGenres)
	router.GET("/movies/most_watched", server.listTheMostWatchedMovies)
	router.GET("/movies/latest", server.listTheLatestReleasedMovies)
	router.GET("/comments", server.listComments)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/movies", server.createMovie)
	authRoutes.PUT("/movies/:id", server.updateMovie)
	authRoutes.POST("/comments", server.createComment)
	authRoutes.PUT("/comments", server.updateComment)
	authRoutes.DELETE("/comments", server.deleteComment)

	server.router = router
}

//Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
