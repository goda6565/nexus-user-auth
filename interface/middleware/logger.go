package middleware

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"github.com/goda6565/nexus-user-auth/pkg/logger"
)

// GinZap は、HTTPリクエストのログ出力ミドルウェアを返します。
func GinZap() gin.HandlerFunc {
	return ginzap.Ginzap(logger.ZapLogger, time.RFC3339, true)
}

// RecoveryWithZap は、パニックからの回復とエラーログ出力のミドルウェアを返します。
func RecoveryWithZap() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger.ZapLogger, true)
}
