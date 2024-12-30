package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeOutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				c.AbortWithStatus(http.StatusGatewayTimeout)
			}
			cancel()
		}()
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
