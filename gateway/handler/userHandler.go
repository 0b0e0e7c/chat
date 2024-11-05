package handler

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"chatting/service/userService/pb/user"

	"github.com/gin-gonic/gin"
)

// Register handles the registration of a new user.
//
//	@Summary		Register a new user
//	@Description	Register a new user with a username and password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			json_object	body		user.RegisterRequest	true	"register request"
//	@Success		200			{object}	_BaseResponse{data=user.RegisterResponse}
//	@Failure		500			{object}	_ResponseError
//	@Router			/user/register [post]
func (l *Logic) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := l.UserRpc().Register(
		context.Background(), &user.RegisterRequest{
			Username: req.Username,
			Password: req.Password,
		},
	)

	HandleRespondsWithError(c, resp, err)
}

// Login handles the logging in of a user.
//
//	@Summary		Login a user
//	@Description	Login a user with a username and password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			login	body		user.LoginRequest	true	"login request"
//	@Success		200		{object}	_BaseResponse{data=user.LoginResponse}
//	@Failure		500		{object}	_ResponseError
//	@Router			/user/login [post]
func (l *Logic) Login(c *gin.Context) {
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := l.UserRpc().Login(
		context.Background(), &req,
	)

	logx.Info("login response: ", resp)

	HandleRespondsWithError(c, resp, err)
}

// Logout handles the logging out of a user.
//
//	@Summary		Logout a user
//	@Description	Logout a user by invalidating their token
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			token	body		string	true	"Token"
//	@Success		200		{object}	_BaseResponse{data=user.LogoutResponse}
//	@Failure		500		{object}	_ResponseError
//	@Router			/user/logout [post]
func (l *Logic) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := l.UserRpc().Logout(
		context.Background(), &user.LogoutRequest{
			UserId: userID,
			Token:  req.Token})

	HandleRespondsWithError(c, resp, err)
}

func (l *Logic) UserIdExist(c *gin.Context, userID []int64) (error, bool) {
	_, err := l.userRpc.CheckUserExist(context.Background(), &user.CheckUserExistRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"check user(s) exist error": err.Error()})
		return nil, false
	}
	return nil, true
}
