package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"chatting/gateway/handler"
	"chatting/service/userService/pb/user"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	userRPCClient user.UserServiceClient
	authHeader    string
	tokenPrefix   string
}

func NewJWTMiddleware() gin.HandlerFunc {
	return (&JWTMiddleware{
		userRPCClient: handler.UserRpcClient,
		authHeader:    "Authorization",
		tokenPrefix:   "Bearer ",
	}).Handle()
}

func (m *JWTMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取 JWT Token
		token, err := m.extractToken(c)
		if err != nil {
			m.respondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		// 调用 RPC 方法进行 Token 验证
		resp, err := m.userRPCClient.ValidateJWT(context.Background(), &user.ValidateRequest{Token: token})
		if err != nil || !resp.Valid {
			m.respondWithError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", resp.UserId)
		c.Next()
	}
}

// extractToken 提取 Bearer Token
func (m *JWTMiddleware) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(m.authHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// 检查是否包含 Bearer 前缀
	if !strings.HasPrefix(authHeader, m.tokenPrefix) {
		return "", errors.New("bearer token is required")
	}

	// 返回去除前缀后的 Token
	return strings.TrimPrefix(authHeader, m.tokenPrefix), nil
}

func (m *JWTMiddleware) respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
	c.Abort()
}
