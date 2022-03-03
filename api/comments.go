package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	db "phantom/db/mongo"
	"phantom/token"
)

var errUnAuthorizedUser = errors.New("Unauthorized User")
var errCommentNotFound = errors.New("Comment is not found")

type createCommentRequest struct {
	Email   string `json:"email" binding:"email|max=0"`
	MovieID string `json:"movie_id" binding:"required,hexadecimal,min=24"`
	Text    string `json:"text" binding:"required,min=1"`
}

func (server *Server) createComment(ctx *gin.Context) {
	var req createCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	objectId, err := primitive.ObjectIDFromHex(req.MovieID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.AddCommentParams{
		Name:    authPayload.Username,
		Email:   req.Email,
		MovieID: objectId,
		Text:    req.Text,
	}
	id, err := server.store.AddComment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id.Hex()})
}

type listCommentsRequest struct {
	MovieID  string `form:"movie_id" binding:"required,hexadecimal,min=24"`
	PageSize int64  `form:"s" binding:"required,min=1,max=20"`
	PageId   int64  `form:"p" binding:"required,min=1"`
}
type ListCommentsResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
}

func (server *Server) listComments(ctx *gin.Context) {
	var req listCommentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	objectId, err := primitive.ObjectIDFromHex(req.MovieID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.GetCommentsParams{
		MovieID: objectId,
		Skip:    req.PageSize * (req.PageId - 1),
		Limit:   req.PageSize,
	}
	comments, err := server.store.GetCommentsByMovieID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var rsp []ListCommentsResponse
	for _, comment := range comments {
		rsp = append(rsp, ListCommentsResponse{
			Id:   comment.ID.Hex(),
			Name: comment.Name,
			Text: comment.Text,
		})
	}
	ctx.JSON(http.StatusOK, rsp)
}

type updateCommentRequest struct {
	Id   string `json:"id" binding:"required,hexadecimal,min=24"`
	Text string `json:"text" binding:"required,min=1"`
}

func (server *Server) updateComment(ctx *gin.Context) {
	var req updateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.Comments{
		ID:   objectId,
		Name: authPayload.Username,
		Text: req.Text,
	}
	_, err = server.store.UpdateComment(ctx, arg)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			ctx.JSON(http.StatusNotFound, errorResponse(errCommentNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"updated": "ok"})
}

type commentIdRequest struct {
	Id string `json:"id" binding:"required,hexadecimal,min=24"`
}

func (server *Server) deleteComment(ctx *gin.Context) {
	var req commentIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	_, err = server.store.DeleteComment(ctx, objectId, authPayload.Username)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			ctx.JSON(http.StatusNotFound, errorResponse(errCommentNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted": "OK"})
}
