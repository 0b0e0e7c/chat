package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"

	"chatting/service/friendService/pb/friend"
	"chatting/service/userService/pb/user"
)

// PendingFriends handles the pending friend requests.
//
//	@Summary		Send Pending friend requests
//	@Description	Send Pending friend requests
//	@Tags			Friend
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		friend.AddFriendRequest	true	"Friend ID"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/pending [post]
func (l *Logic) PendingFriends(c *gin.Context) {
	var req friend.AddFriendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		return
	}

	_, ok := l.UserIdExist(c, []int64{userID, req.FriendId})
	if !ok {
		return
	}

	logx.Info("userID: ", userID)

	resp, err := l.FriendRpc().PendingRequest(context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// Confirm handles the confirmation of a friend request.
//
//	@Summary		Add a friend
//	@Description	Add a friend by their user ID
//	@Tags			Friend
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		friend.AddFriendRequest	true	"Friend ID"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/confirm [post]
func (l *Logic) Confirm(c *gin.Context) {
	var req friend.AddFriendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	_, ok := l.UserIdExist(c, []int64{userID, req.FriendId})
	if !ok {
		return
	}

	logx.Info("userID: ", userID)

	resp, err := l.FriendRpc().ConfirmRequest(context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

type friends struct {
	UserId   int64  `json:"friend_id" example:"1"`
	Username string `json:"friend_username" example:"user1"`
}

// GetFriends handles the getting of a user's friends.
//
//	@Summary		Get a user's friends
//	@Description	Get a user's friends
//	@Tags			Friend
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	_BaseResponse{data=friends}
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/friends [get]
func (l *Logic) GetFriends(c *gin.Context) {
	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	_, ok := l.UserIdExist(c, []int64{userID})
	if !ok {
		return
	}

	// Fetch friends via the friend RPC
	resp, err := l.FriendRpc().GetFriends(
		context.Background(), &friend.GetFriendsRequest{
			UserId: userID,
		})
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	// Extract friend IDs
	var friendIDs []int64
	for _, g := range resp.Friends {
		friendIDs = append(friendIDs, g.UserId)
	}

	friendsList, err := l.mapUsernames(context.Background(), friendIDs)
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	// Respond with the mapped friends list
	HandleRespondsWithError(c, friendsList, nil)
}

// GetPendingFriends handles the getting of a user's pending friends.
//
//	@Summary		Get a user's pending friends
//	@Description	Get a user's pending friends
//	@Tags			Friend
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	_BaseResponse{data=friends}
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/pending [get]
func (l *Logic) GetPendingFriends(c *gin.Context) {
	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	_, ok := l.UserIdExist(c, []int64{userID})
	if !ok {
		return
	}

	// Fetch pending friend requests via the friend RPC
	resp, err := l.FriendRpc().GetPendingRequests(
		context.Background(), &friend.GetFriendsRequest{
			UserId: userID,
		})
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	// Extract pending friend IDs
	var pendingFriendIDs []int64
	for _, pr := range resp.Friends { // Assuming resp.PendingFriends exists
		pendingFriendIDs = append(pendingFriendIDs, pr.UserId)
	}

	// Map usernames using the helper function
	pendingFriendsList, err := l.mapUsernames(context.Background(), pendingFriendIDs)
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	// Respond with the mapped pending friends list
	HandleRespondsWithError(c, pendingFriendsList, nil)
}

func (l *Logic) mapUsernames(ctx context.Context, userIDs []int64) ([]friends, error) {
	if len(userIDs) == 0 {
		return []friends{}, nil
	}

	// Fetch usernames using the user RPC
	nameResp, err := l.userRpc.GetUserNameById(ctx, &user.GetUserNameByIdRequest{UserId: userIDs})
	if err != nil {
		return nil, err
	}

	mapping := nameResp.UserNames

	// Map user IDs to usernames
	var friendsList []friends
	for _, id := range userIDs {
		username, exists := mapping[id]
		if !exists {
			username = "Unknown"
		}
		friendsList = append(
			friendsList, friends{
				UserId:   id,
				Username: username,
			})
	}
	return friendsList, nil
}

// DeleteFriend handles the deleting of a friend.
//
//	@Summary		Delete a friend
//	@Description	Delete a friend by their user ID
//	@Tags			Friend
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		friend.DeleteFriendRequest	true	"Friend ID"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/delete [post]
func (l *Logic) DeleteFriend(c *gin.Context) {
	var req friend.DeleteFriendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		return
	}
	req.UserId = userID

	err, ok := l.UserIdExist(c, []int64{userID, req.FriendId})
	if !ok {
		return
	}

	resp, err := l.FriendRpc().DeleteFriend(context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// CreateFriendGroup handles the creation of a friend group.
//
//	@Summary		Create a friend group
//	@Description	Create a friend group
//	@Tags			Friend/manage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		friend.CreateFriendGroupRequest	true	"CreateFriendGroupRequest request"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/manage/createGroup [post]
func (l *Logic) CreateFriendGroup(c *gin.Context) {
	var req friend.CreateFriendGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		return
	}
	req.UserId = userID

	resp, err := l.FriendRpc().CreateFriendGroup(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// AddToGroup handles the adding of a friend to a friend group.
//
//	@Summary		Add a friend to a group
//	@Description	Add a friend to a group
//	@Tags			Friend/manage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		friend.AddFriendToGroupRequest	true	"AddToGroupRequest request"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/manage/addToGroup [post]
func (l *Logic) AddToGroup(c *gin.Context) {
	var req friend.AddFriendToGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		return
	}
	req.UserId = userID

	err, ok := l.UserIdExist(c, []int64{userID, req.FriendId})
	if !ok {
		return
	}

	resp, err := l.FriendRpc().AddFriendToGroup(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// DeleteFromGroup handles the deleting of a friend from a friend group.
//
//	@Summary		Delete a friend from a group
//	@Description	Delete a friend from a group
//	@Tags			Friend/manage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		friend.AddFriendToGroupRequest	true	"DeleteFromGroupRequest request"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/manage/deleteFromGroup [post]
func (l *Logic) DeleteFromGroup(c *gin.Context) {
	var req friend.AddFriendToGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		return
	}

	err, ok := l.UserIdExist(c, []int64{userID, req.FriendId})
	if !ok {
		return
	}

	resp, err := l.FriendRpc().DeleteFriendFromGroup(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

type friendGroup struct {
	GroupID   int64     `json:"group_id"`
	GroupName string    `json:"group_name"`
	Friends   []friends `json:"friends"`
}

// GetFriendGroup handles the getting of a user's friend groups.
//
//	@Summary		Get a user's friend groups
//	@Description	Get a user's friend groups
//	@Tags			Friend/manage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	_BaseResponse{data=[]friendGroup}
//	@Failure		500			{object}	_ResponseError
//	@Router			/friend/manage/friendGroup [get]
func (l *Logic) GetFriendGroup(c *gin.Context) {
	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		return
	}

	logx.Info("get friend group for userID: ", userID)

	resp, err := l.FriendRpc().GetFriendGroup(context.Background(), &friend.GetFriendGroupRequest{UserId: userID})
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	var friendIDs []int64
	for _, g := range resp.Groups {
		for _, f := range g.Friends {
			friendIDs = append(friendIDs, f.UserId)
		}
	}

	nameResp, err := l.userRpc.GetUserNameById(context.Background(), &user.GetUserNameByIdRequest{UserId: friendIDs})
	if err != nil {
		HandleRespondsWithError(c, nil, err)
		return
	}

	mapping := nameResp.UserNames

	var friendGroups []friendGroup
	for _, g := range resp.Groups {
		var friendsList []friends
		for _, f := range g.Friends {
			friendsList = append(
				friendsList, friends{
					UserId:   f.UserId,
					Username: mapping[f.UserId],
				})
		}

		friendGroups = append(
			friendGroups, friendGroup{
				GroupID:   g.GroupId,
				GroupName: g.GroupName,
				Friends:   friendsList,
			})
	}

	HandleRespondsWithError(c, friendGroups, nil)
}

func (l *Logic) UserAreFriend(c *gin.Context, userID, friendID int64) error {
	_, err := l.FriendRpc().CheckIsFriend(
		context.Background(), &friend.CheckIsFriendRequest{UserId: userID, FriendId: friendID})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"logic error": err.Error()})
		return err
	}
	return nil
}
