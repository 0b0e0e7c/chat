package logic

import (
	"context"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsLogic {
	return &GetFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsLogic) GetFriends(in *friend.GetFriendsRequest) (*friend.GetFriendsResponse, error) {
	friendDao := l.svcCtx.FriendDao
	if in.GetUserId() == 0 {
		return nil, ErrUserIdRequired
	}

	// 获取好友列表
	var friends []model.Friend
	friends, err := friendDao.GetFriendsWithStatus(in.UserId, model.Added)
	if err != nil {
		return nil, err
	}

	// 构建响应
	resp := &friend.GetFriendsResponse{}

	for _, f := range friends {
		if f.FriendID == in.UserId {
			f.FriendID, f.UserID = f.UserID, f.FriendID
		}

		resp.Friends = append(
			resp.Friends, &friend.FriendItem{
				UserId: f.FriendID,
			})

	}

	return resp, nil
}
