package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"

	g "chatting/service/groupMessageService/pb/groupMessage"
)

type GroupMessageHandler struct {
	logic          *Logic
	wsGroupClients map[int64]map[int64]*websocket.Conn // map[groupID]map[userID]*websocket.Conn
	wsMutex        sync.RWMutex
	upgrader       websocket.Upgrader
}

func NewGroupMessageHandler() *GroupMessageHandler {
	return &GroupMessageHandler{
		logic:          GlobalLogic,
		wsGroupClients: make(map[int64]map[int64]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

// CreateMsgGroup handles the creation of a message group.
//
//	@Summary		Create a message group
//	@Description	Create a group with a list of members
//	@Tags			GroupMessage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		g.CreateGroupRequest	true	"Group name and member IDs"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/groupMessage/create [post]
func (h *GroupMessageHandler) CreateMsgGroup(c *gin.Context) {
	var req g.CreateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	req.UserId = userID

	resp, err := h.logic.GroupMsgRpc().CreateGroup(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// AddGroupMember handles the adding of a member to a group.
//
//	@Summary		Add a member to a group
//	@Description	Add a member to a group
//	@Tags			GroupMessage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		g.AddGroupMemberRequest	true	"Group ID and member ID"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/groupMessage/addGroupMember [post]
func (h *GroupMessageHandler) AddGroupMember(c *gin.Context) {
	var req g.AddGroupMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	req.UserId = userID

	resp, err := h.logic.GroupMsgRpc().AddGroupMember(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// SendMsg handles the sending of a group message.
//
//	@Summary		Send a group message
//	@Description	Send a message to a group
//	@Tags			GroupMessage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		g.SendGroupMessageRequest	true	"Group ID and content"
//	@Success		200			{object}	_BaseResponse
//	@Failure		500			{object}	_ResponseError
//	@Router			/groupMessage/send [post]
func (h *GroupMessageHandler) SendMsg(c *gin.Context) {
	var req g.SendGroupMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	req.SenderId = userID

	resp, err := h.logic.GroupMsgRpc().SendGroupMessage(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

// GetMsg handles the getting of group messages.
//
//	@Summary		Get group messages
//	@Description	Get messages from a group
//	@Tags			GroupMessage
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			json_object	body		g.GetGroupMessagesRequest	true	"Group ID, limit, and offset"
//	@Success		200			{object}	_BaseResponse{data=[]g.GroupMessage}
//	@Failure		500			{object}	_ResponseError
//	@Router			/groupMessage [get]
func (h *GroupMessageHandler) GetMsg(c *gin.Context) {
	var req g.GetGroupMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserIDFromGinCtx(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	req.UserId = userID

	resp, err := h.logic.GroupMsgRpc().GetGroupMessages(
		context.Background(), &req)

	HandleRespondsWithError(c, resp, err)
}

func (h *GroupMessageHandler) WebSocketHandler(c *gin.Context) {
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

	h.removeClientFromGroup(userID, -1) // 清理所有群组
}

func (h *GroupMessageHandler) handleWebSocketMessage(userID int64, msg map[string]any, conn *websocket.Conn) {
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
	case "join_group":
		groupID, ok := data["group_id"].(int64)
		if !ok {
			logx.Error("Invalid group_id")
			return
		}
		h.handleJoinGroup(userID, int64(groupID), conn)
	case "send_group_message":
		groupID, ok := data["group_id"].(int64)
		if !ok {
			logx.Error("Invalid group_id")
			return
		}
		content, ok := data["content"].(string)
		if !ok {
			logx.Error("Invalid content")
			return
		}
		h.handleSendGroupMessage(userID, int64(groupID), content)
	default:
		logx.Errorf("Unknown message type: %s", msgType)
	}
}

func (h *GroupMessageHandler) handleJoinGroup(userID, groupID int64, conn *websocket.Conn) {
	h.wsMutex.Lock()
	defer h.wsMutex.Unlock()

	if h.wsGroupClients[groupID] == nil {
		h.wsGroupClients[groupID] = make(map[int64]*websocket.Conn)
	}

	h.wsGroupClients[groupID][userID] = conn
	logx.Infof("User %d joined group %d", userID, groupID)
}

func (h *GroupMessageHandler) handleSendGroupMessage(userID, groupID int64, content string) {
	resp, err := h.logic.GroupMsgRpc().SendGroupMessage(
		context.Background(), &g.SendGroupMessageRequest{
			SenderId: userID,
			GroupId:  groupID,
			Content:  content,
		})
	if err != nil {
		logx.Error("Error sending group message via gRPC:", err)
		return
	}

	h.wsMutex.RLock()
	defer h.wsMutex.RUnlock()

	if clients, ok := h.wsGroupClients[groupID]; ok {
		for _, conn := range clients {
			err := conn.WriteJSON(
				gin.H{
					"msg_id":    resp.MsgId,
					"sender_id": userID,
					"group_id":  groupID,
					"content":   content,
					"timestamp": resp.Timestamp,
				})
			if err != nil {
				logx.Error("Error sending message to WebSocket:", err)
			}
		}
	}
}

func (h *GroupMessageHandler) removeClientFromGroup(userID, groupID int64) {
	h.wsMutex.Lock()
	defer h.wsMutex.Unlock()

	if groupID == -1 { // 清理所有群组
		for gid, clients := range h.wsGroupClients {
			delete(clients, userID)
			if len(clients) == 0 {
				delete(h.wsGroupClients, gid)
			}
		}
	} else {
		if clients, ok := h.wsGroupClients[groupID]; ok {
			delete(clients, userID)
			if len(clients) == 0 {
				delete(h.wsGroupClients, groupID)
			}
		}
	}
}
