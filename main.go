package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/goda6565/nexus-user-auth/infrastructure/database"
	"github.com/goda6565/nexus-user-auth/infrastructure/web"
	"github.com/goda6565/nexus-user-auth/pkg/logger"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

func main() {
	// 開発環境の場合、.env.development ファイルを読み込む
	env := utils.GetEnvDefault("ENV", "development")
	if env == "development" {
		err := godotenv.Load(".env.development")
		if err != nil {
			logger.Error("Error loading .env.development file")
		}
	}

	db, err := database.NewDBInstance(database.InstancePostgres)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	server, err := web.NewServer(db)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	// サーバーを非同期で起動する。
	// ※ server.Start() はブロックする処理のため、ゴルーチン内で実行し、
	//    メインゴルーチンが終了シグナル待ちなどの処理を実行できるようにする。
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// エラーが発生した場合、かつそれがサーバー正常終了時のエラーでなければログ出力後、プログラム終了
			logger.Fatal(err.Error())
		}
	}()

	// 終了シグナル（SIGINT や SIGTERM）を受信するためのチャネルを作成する。
	// これにより、Ctrl+C などのシグナル受信でプログラムを終了できるようにする。
	quit := make(chan os.Signal, 1)
	// SIGINT と SIGTERM のシグナルを quit チャネルで受け取るように設定する
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// シグナルが来るまでここで待機する
	<-quit
	log.Println("Shutting down server...")

	// プログラム終了前に、ログバッファに残ったログを全てフラッシュ（出力）する
	defer logger.Sync()

	// 優雅なサーバーのシャットダウン処理のため、タイムアウト付きのコンテキストを作成する。
	// この例では、シャットダウン処理に最大5秒間待機する。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// サーバーのシャットダウン処理を実行する。
	// 現在処理中のリクエストがある場合、処理完了まで待機しながら停止する（Graceful Shutdown）。
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}

	// シャットダウン処理が完了するまで待機する。
	// ctx.Done() が閉じられると、シャットダウンが完了したと判断できる。
	<-ctx.Done()
}
