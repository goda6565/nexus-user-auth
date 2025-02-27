package logger

import (
	"os"

	"go.uber.org/zap"
)

var (
	ZapLogger        *zap.Logger
	zapSugaredLogger *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionConfig() // 本番環境用の設定を読み込む(info以上のログが出力される)
	logFile := os.Getenv("APP_LOG_FILE")
	if logFile != "" {
		cfg.OutputPaths = []string{"stderr", logFile} // 環境変数APP_LOG_FILEが設定されている場合は、標準出力とファイルに出力
	}

	ZapLogger = zap.Must(cfg.Build()) // Loggerを生成(エラーがあればpanic)
	if os.Getenv("APP_ENV") == "development" {
		ZapLogger = zap.Must(zap.NewDevelopment()) // 開発環境の場合は、開発用の設定を読み込む(debug以上のログが出力される)
	}
	zapSugaredLogger = ZapLogger.Sugar() // SugaredLoggerを生成
}

func Sync() {
	err := zapSugaredLogger.Sync() // ログに書き込まれている内容をフラッシュして、バッファリングされたエントリを出力
	if err != nil {
		zap.Error(err)
	}
}

func Info(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Infow(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Panicw(msg, keysAndValues...)
}
