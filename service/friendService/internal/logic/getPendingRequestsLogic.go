package logic

import (
	"context"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPendingRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPendingRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPendingRequestsLogic {
	return &GetPendingRequestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPendingRequestsLogic) GetPendingRequests(in *friend.GetFriendsRequest) (*friend.GetFriendsResponse, error) {
	if in.GetUserId() == 0 {
		return nil, ErrUserIdRequired
	}

	friendDao := l.svcCtx.FriendDao

	// 获取好友列表
	friends, err := friendDao.GetFriendsWithStatus(in.GetUserId(), model.Pending)
	if err != nil {
		return nil, err
	}

	// 构建响应
	resp := &friend.GetFriendsResponse{}

	for _, f := range friends {
		if f.InitiatorID == in.GetUserId() {
			continue
		}

		if f.FriendID == in.GetUserId() {
			f.FriendID, f.UserID = f.UserID, f.FriendID
		}

		resp.Friends = append(
			resp.Friends, &friend.FriendItem{
				UserId: f.FriendID,
			})

	}
	return resp, nil
}
