package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func GetUserIDFromGinCtx(c *gin.Context) (int64, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user not found in gin context")
	}
	return userID.(int64), nil
}

func handleError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
		return
	}
	c.JSON(
		http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  st.Message(),
		})
}

func HandleRespondsWithError(c *gin.Context, resp any, err error) {
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": resp})
}
