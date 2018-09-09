package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check for incoming header, use it if exists
		requestID := c.Request.Header.Get("X-Request-Id")

		// create request id with uuid4
		if requestID == "" {
			u4 := uuid.NewV4()
			requestID = u4.String()
		}

		// expose it for use in the application
		c.Set("X-Request-Id", requestID)

		// set x-request-id header
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}
