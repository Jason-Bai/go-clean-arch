package middleware

import (
	"github.com/Jason-Bai/go-clean-arch/handler"
	"github.com/Jason-Bai/go-clean-arch/pkg/errno"
	"github.com/Jason-Bai/go-clean-arch/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
