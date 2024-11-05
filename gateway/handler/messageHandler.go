package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"

	"chatting/service/messageService/pb/message"
)

type MessageHandler struct {
	logic     *Logic
	wsClients map[int64]*websocket.Conn
	wsMutex   sync.RWMutex
	upgrader  websocket.Upgrader
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		logic:     GlobalLogic,
		wsClients: make(map[int64]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

// SendMsg handles the sending of a message.
//
//	@Summary		Send a message
//	@Description	Send a message to a friend
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		message.SendMessageRequest	true	"Receiver ID and content"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/message/send [post]
func (h *MessageHandler) SendMsg(c *gin.Context) {
	var req message.SendMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	req.SenderId = userID

	if err = h.logic.UserAreFriend(c, userID, req.ReceiverId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "users are not friend"})
		return
	}

	resp, err := h.logic.MsgRpc().SendMessage(context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// GetMsg handles the fetching of messages.
//
//	@Summary		Get messages
//	@Description	Get messages between a user and a friend
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		message.GetMessagesRequest	true	"Peer ID, offset and  limit"
//	@Success		200			{object}	_BaseResponse{data=message.GetMessagesResponse}
//	@Failure		500			{object}	_ResponseError
//	@Router			/message [get]
func (h *MessageHandler) GetMsg(c *gin.Context) {
	var req message.GetMessagesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	req.UserId = userID

	if err = h.logic.UserAreFriend(c, userID, req.PeerId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "users are not friend"})
		return
	}

	resp, err := h.logic.MsgRpc().GetMessages(context.Background(), &req)

	HandleRespondsWithError(c, resp, err)

}

func (h *MessageHandler) WebSocketHandler(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logx.Error("Error upgrading to WebSocket:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logx.Error("Error closing WebSocket connection:", err)
		}
	}(conn)

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		logx.Error("User ID not found in context")
		return
	}

	h.wsMutex.Lock()
	h.wsClients[userID] = conn
	h.wsMutex.Unlock()

	for {
		_, rMsg, err := conn.ReadMessage()
		if err != nil {
			logx.Error("Error reading message from WebSocket:", err)
			break
		}

		var msg map[string]any
		if err := json.Unmarshal(rMsg, &msg); err != nil {
			logx.Error("Error parsing JSON message:", err)
			continue
		}

		h.handleWebSocketMessage(userID, msg, conn)
	}

	h.wsMutex.Lock()
	delete(h.wsClients, userID)
	h.wsMutex.Unlock()
}

func (h *MessageHandler) handleWebSocketMessage(userID int64, msg map[string]any, conn *websocket.Conn) {
	msgType, ok := msg["type"].(string)
	if !ok {
		logx.Error("Invalid message type")
		return
	}

	data, ok := msg["data"].(map[string]any)
	if !ok {
		logx.Error("Invalid message data")
		return
	}

	switch msgType {
	case "start_chat":
		peerID, ok := data["peer_id"].(float64)
		if !ok {
			logx.Error("Invalid peer_id")
			return
		}
		h.handleStartChat(userID, int64(peerID), conn)
	case "send_message":
		peerID, ok := data["peer_id"].(float64)
		if !ok {
			logx.Error("Invalid peer_id")
			return
		}
		content, ok := data["content"].(string)
		if !ok {
			logx.Error("Invalid content")
			return
		}
		h.handleSendMessage(userID, int64(peerID), content, conn)
	default:
		logx.Errorf("Unknown message type: %s", msgType)
	}
}

func (h *MessageHandler) handleStartChat(userID, peerID int64, conn *websocket.Conn) {
	resp, err := h.logic.MsgRpc().GetMessages(
		context.Background(), &message.GetMessagesRequest{
			UserId: userID,
			PeerId: peerID,
			Limit:  10,
			Offset: -1,
		})
	if err != nil {
		logx.Error("Error fetching messages:", err)
		return
	}

	for _, msg := range resp.Messages {
		if err := conn.WriteJSON(msg); err != nil {
			logx.Error("Error sending message to WebSocket:", err)
			return
		}
	}
}

func (h *MessageHandler) handleSendMessage(userID, peerID int64, content string, conn *websocket.Conn) {
	resp, err := h.logic.MsgRpc().SendMessage(
		context.Background(), &message.SendMessageRequest{
			SenderId:   userID,
			ReceiverId: peerID,
			Content:    content,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			logx.Error("Error sending message via gRPC:", err)
			if err := conn.WriteJSON(gin.H{"status": "failed", "msg": err.Error()}); err != nil {
				logx.Error("Error sending message to WebSocket:", err)
			}
		} else {
			logx.Error("Error sending message via gRPC:", st.Message())
			if err := conn.WriteJSON(gin.H{"status": "failed", "msg": st.Message()}); err != nil {
				logx.Error("Error sending message to WebSocket:", err)
			}
		}
		return
	}

	if err := conn.WriteJSON(gin.H{"status": "success", "msg": "Message sent successfully"}); err != nil {
		logx.Error("Error sending message to WebSocket:", err)
	}

	h.notifyReceiver(resp.MsgId, peerID, userID, resp.Timestamp, content)
}

func (h *MessageHandler) notifyReceiver(msgID, receiverID, senderID, timestamp int64, content string) {
	h.wsMutex.RLock()
	defer h.wsMutex.RUnlock()

	if receiverConn, ok := h.wsClients[receiverID]; ok {
		err := receiverConn.WriteJSON(
			gin.H{
				"msg_id":      msgID,
				"sender_id":   senderID,
				"receiver_id": receiverID,
				"content":     content,
				"timestamp":   timestamp,
			})
		if err != nil {
			logx.Error("Error sending message to receiver via WebSocket:", err)
		}
	}
}
