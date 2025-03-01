package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/goda6565/nexus-user-auth/interface/gen"
	"github.com/goda6565/nexus-user-auth/interface/keys"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

func AuthMiddleware(path string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 指定されたパスと完全一致しなければ認証処理をスキップ
		if c.Request.URL.Path != path {
			c.Next()
			return
		}

		// Authorization ヘッダーの取得
		authHeader := c.Request.Header.Get("Authorization")

		// Bearer プレフィックスを除去
		if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = authHeader[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gen.ErrorResponse{
				Message: "Invalid token",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		// トークン検証
		claims, err := utils.ValidateToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gen.ErrorResponse{
				Message: "Invalid token",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		// Gin の Context にユーザーIDをセット
		c.Set("validated_uid", claims.ObjID)

		// リクエストのコンテキストも更新
		newCtx := context.WithValue(c.Request.Context(), keys.ValidatedUIDKey, claims.ObjID)
		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	}
}
