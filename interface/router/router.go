package router

import (
	"context"
	"encoding/json"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"github.com/goda6565/nexus-user-auth/application/service/user/authentication"
	"github.com/goda6565/nexus-user-auth/application/service/user/profile"
	"github.com/goda6565/nexus-user-auth/application/service/user/registration"
	"github.com/goda6565/nexus-user-auth/infrastructure/database/repository"
	"github.com/goda6565/nexus-user-auth/interface/gen"
	"github.com/goda6565/nexus-user-auth/interface/handler"
	"github.com/goda6565/nexus-user-auth/interface/handler/user"
	"github.com/goda6565/nexus-user-auth/interface/middleware"
	"github.com/goda6565/nexus-user-auth/pkg/logger"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

type ServerInterfaceImpl struct {
	*user.UserRegistrationHandler
	*user.UserAuthenticationHandler
	*user.UserProfileHandler
}

// swagger設定
func setUpSwagger(router *gin.Engine) (*openapi3.T, error) {
	// OpenAPI (Swagger) 定義を取得
	swagger, err := gen.GetSwagger()
	if err != nil {
		return nil, err
	}

	// 環境変数 ENV を取得（デフォルトは "development"）
	env := utils.GetEnvDefault("ENV", "development")
	if env == "development" {
		// Swagger 定義を JSON に変換
		swaggerJson, _ := json.Marshal(swagger)
		// swag 用の Spec オブジェクトを作成
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",           // インスタンス名
			SwaggerTemplate:  string(swaggerJson), // Swagger のテンプレート（JSON 文字列）
		}
		// swag ライブラリに登録
		swag.Register(SwaggerInfo.InfoInstanceName, SwaggerInfo)
		// /swagger/*any で Swagger UI を提供
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	// Swagger 定義を返す
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	// Gin Engine を作成
	router := gin.New()

	// CORS 設定
	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	// Swagger の設定（開発環境の場合は Swagger UI を有効化）
	swagger, err := setUpSwagger(router)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// ミドルウェアの設定
	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	// health check
	router.GET("/health", handler.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(10 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			// OapiRequestValidatorWithOptions を利用して、認証関数付きのバリデーションミドルウェアを作成
			// TODO: 認証関数を実装する（ミドルウェアで認証を完了してユーザ情報を取得）
			v1.Use(ginMiddleware.OapiRequestValidatorWithOptions(swagger, &ginMiddleware.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
						return nil
					},
				},
			}))
			userRepositoryImpl := repository.NewUserRepository(db)
			// ユーザー登録サービスを作成
			userRegistrationService := registration.NewUserRegistrationService(userRepositoryImpl)
			// ユーザー登録ハンドラを作成
			userRegistrationHandler := user.NewUserRegistrationHandler(userRegistrationService)
			// ユーザー認証サービスを作成
			userAuthenticationService := authentication.NewUserAuthenticationService(userRepositoryImpl)
			// ユーザー認証ハンドラを作成
			userAuthenticationHandler := user.NewUserAuthenticationHandler(userAuthenticationService)
			// ユーザープロフィールサービスを作成
			userProfileService := profile.NewUserProfileService(userRepositoryImpl)
			// ユーザープロフィールハンドラを作成
			userProfileHandler := user.NewUserProfileHandler(userProfileService)

			// ハンドラをまとめる
			serverInterface := &ServerInterfaceImpl{
				UserRegistrationHandler:   userRegistrationHandler,
				UserAuthenticationHandler: userAuthenticationHandler,
				UserProfileHandler:        userProfileHandler,
			}

			gen.RegisterHandlers(v1, serverInterface)
		}
	}
	// ルーターを返す
	return router, nil
}
