package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		logx.WithContext(ctx).Info(path,
			fmt.Sprintf(" result_type: %v", "request"),
			fmt.Sprintf(" method: %v", c.Request.Method),
			fmt.Sprintf(" path: %v", path),
			fmt.Sprintf(" query: %v", query),
			fmt.Sprintf(" ip: %v", c.ClientIP()),
		)

		c.Next()

		var errInf string
		if len(c.Errors) > 0 {
			errInf = c.Errors.ByType(gin.ErrorTypePrivate).String()
		} else {
			errInf = ""
		}

		cost := time.Since(start)
		logx.WithContext(ctx).Info(path,
			fmt.Sprintf(" result_type: %v", "response"),
			fmt.Sprintf(" status: %v", c.Writer.Status()),
			fmt.Sprintf(" method: %v", c.Request.Method),
			fmt.Sprintf(" path: %v", path),
			fmt.Sprintf(" query: %v", query),
			fmt.Sprintf(" ip: %v", c.ClientIP()),
			fmt.Sprintf(" user-agent: %v", c.Request.UserAgent()),
			fmt.Sprintf(" errors: %v", errInf),
			fmt.Sprintf(" cost: %.1fms", float32(cost)/float32(time.Millisecond)),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logx.WithContext(ctx).Error(c.Request.URL.Path,
						fmt.Sprintf("error: %+v", err),
						fmt.Sprintf("request: %v", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logx.WithContext(ctx).Error("[Recovery from panic]",
						fmt.Sprintf("error: %+v", err),
						fmt.Sprintf("request: %v", string(httpRequest)),
						fmt.Sprintf("stack: %v", string(debug.Stack())),
					)
				} else {
					logx.WithContext(ctx).Error("[Recovery from panic]",
						fmt.Sprintf("error: %+v", err),
						fmt.Sprintf("request: %v", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
