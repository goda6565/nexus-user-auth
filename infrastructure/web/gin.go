package web

import (
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/goda6565/nexus-user-auth/interface/router"
	"github.com/goda6565/nexus-user-auth/pkg/logger"
)

type GinWebServer struct {
	server *http.Server
}

func (g *GinWebServer) Start() error {
	return g.server.ListenAndServe()
}

func (g *GinWebServer) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}

func NewGinServer(host, port string, corsAllowOrigins []string, db *gorm.DB) (Server, error) {
	// Gin ルーターの初期化
	router, err := router.NewGinRouter(db, corsAllowOrigins)
	if err != nil {
		logger.Error(err.Error(), "host", host, "port", port)
		return nil, err
	}

	// http.Server を生成して、Gin のルーターをハンドラーに設定
	return &GinWebServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", host, port),
			Handler: router,
		},
	}, err
}
