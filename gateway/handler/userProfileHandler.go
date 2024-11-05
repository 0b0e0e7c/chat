package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"

	"chatting/component/common"
	"chatting/service/userService/pb/user"
)

// UploadAvatar handles the uploading of a user's avatar.
//
//	@Summary		Upload a user's avatar
//	@Description	Upload a user's avatar
//	@Tags			User/profile
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			avatar	formData	file	true	"avatar"
//	@Success		200		{object}	_BaseResponse
//	@Failure		500		{object}	_ResponseError
//	@Router			/user/profile/avatar/upload [put]
func (l *Logic) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 打开上传的文件
	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(srcFile multipart.File) {
		err := srcFile.Close()
		if err != nil {
			logx.Errorf("failed to close file: %v", err)
		}
	}(srcFile)

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, srcFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 上传到 MinIO
	bucketName := "avatars"
	objectName := fmt.Sprintf("%v.jpg", common.HashingUserID(userID.(int64)))
	contentType := "image/jpeg"

	if err := l.bucketClient.UploadFileFromBuffer(bucketName, objectName, buf, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(
		http.StatusOK, gin.H{
			"msg": "ok",
			"data": gin.H{
				"message":    "File uploaded successfully",
				"avatar url": fmt.Sprintf("/%s/%s", bucketName, objectName),
			},
		})
}

func (l *Logic) UpdateUserAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	avatar, exists := c.Get("avatar")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar not found"})
		return
	}

	resp, err := l.userRpc.UpdateUserProfile(
		context.Background(), &user.UpdateUserProfileRequest{
			UserId: userID.(int64),
			Avatar: avatar.(string),
		})
	HandleRespondsWithError(c, resp, err)
}

// GetAvatar handles the downloading of a user's avatar.
//
//	@Summary		Download a user's avatar
//	@Description	Download a user's avatar
//	@Tags			User/profile
//	@Produce		jpeg
//	@Security		BearerAuth
//	@Success		200	{string}	string	"avatar"
//	@Failure		500	{object}	_ResponseError
//	@Router			/user/profile/avatar [get]
func (l *Logic) GetAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	avatarName := common.HashingUserID(userID.(int64)) + ".jpg"

	buf, err := l.bucketClient.DownloadFileToBuffer("avatars", avatarName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download avatar", "details": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
}

// UpdateProfile handles the updating of a user's profile.
//
//	@Summary		Update a user's profile
//	@Description	Update a user's profile
//	@Tags			User/profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			avatar	body		user.UpdateUserProfileRequest	true	"updateProfile request"
//	@Success		200		{object}	_BaseResponse
//	@Failure		500		{object}	_ResponseError
//	@Router			/user/profile/update [put]
func (l *Logic) UpdateProfile(c *gin.Context) {
	var req user.UpdateUserProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	resp, err := l.userRpc.UpdateUserProfile(
		context.Background(), &user.UpdateUserProfileRequest{
			UserId:   userID.(int64),
			Avatar:   req.Avatar,
			Nickname: req.Nickname,
			Email:    req.Email,
			Status:   req.Status,
		})
	HandleRespondsWithError(c, resp, err)
}

// GetProfile handles the getting of a user's profile.
//
//	@Summary		Get a user's profile
//	@Description	Get a user's profile
//	@Tags			User/profile
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	_BaseResponse{data=user.GetUserProfileResponse}
//	@Failure		500	{object}	_ResponseError
//	@Router			/user/profile [get]
func (l *Logic) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	resp, err := l.userRpc.GetUserProfile(
		context.Background(), &user.GetUserProfileRequest{
			UserId: userID.(int64),
		})

	HandleRespondsWithError(c, resp, err)
}
