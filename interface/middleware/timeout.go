package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"

	"github.com/goda6565/nexus-user-auth/interface/gen"
)

func TimeoutMiddleware(timeoutDuration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(timeoutDuration),
		timeout.WithHandler(func(c *gin.Context) { // timeoutが発生するまで普通の処理を行う
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) { // timeoutが発生した場合の処理
			c.JSON(http.StatusRequestTimeout, gen.ErrorResponse{
				Message: "Request timeout",
				Code:    http.StatusRequestTimeout,
			})
			c.Abort()
		}),
	)
}
