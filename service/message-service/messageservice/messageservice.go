// Code generated by goctl. DO NOT EDIT.
// Source: message.proto

package messageservice

import (
	"context"

	"github.com/0b0e0e7c/chat/service/message-service/pb/message"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GetMessagesRequest  = message.GetMessagesRequest
	GetMessagesResponse = message.GetMessagesResponse
	Message             = message.Message
	SendMessageRequest  = message.SendMessageRequest
	SendMessageResponse = message.SendMessageResponse

	MessageService interface {
		SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
		GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error)
	}

	defaultMessageService struct {
		cli zrpc.Client
	}
)

func NewMessageService(cli zrpc.Client) MessageService {
	return &defaultMessageService{
		cli: cli,
	}
}

func (m *defaultMessageService) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	client := message.NewMessageServiceClient(m.cli.Conn())
	return client.SendMessage(ctx, in, opts...)
}

func (m *defaultMessageService) GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error) {
	client := message.NewMessageServiceClient(m.cli.Conn())
	return client.GetMessages(ctx, in, opts...)
}
