package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	db "phantom/db/mongo"
	"phantom/util"
)

type registerRequest struct {
	Username string `json:"name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("hashed password err")))
		return
	}

	arg := db.AddUserParams{
		Name:     req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}
	id, err := server.store.AddUser(ctx, arg)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("user already exists")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("add user error")))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": id.Hex()})
}

type loginRequest struct {
	Username string `json:"username" binding:"alphanum|max=0"`
	Email    string `json:"email" binding:"email|max=0"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID       string
	Username string
	Email    string
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:       user.ID.Hex(),
		Username: user.Name,
		Email:    user.Email,
	}
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Username == "" && req.Email == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("not allowed to be empty")))
		return
	}
	var user db.User
	if req.Email != "" {
		u, err := server.store.GetUserByEmail(ctx, req.Email)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user is not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		user = u
	} else {
		u, err := server.store.GetUserByName(ctx, req.Username)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		user = u
	}

	err := util.CheckPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("username or password error")))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Name, server.config.AccessTokenDuration)
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	rsp := loginResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
