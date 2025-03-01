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
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	authenticationService "github.com/goda6565/nexus-user-auth/application/service/user/authentication"
	profileService "github.com/goda6565/nexus-user-auth/application/service/user/profile"
	registrationService "github.com/goda6565/nexus-user-auth/application/service/user/registration"
	"github.com/goda6565/nexus-user-auth/infrastructure/database/repository"
	"github.com/goda6565/nexus-user-auth/interface/gen"
	"github.com/goda6565/nexus-user-auth/interface/handler"
	authenticationHandler "github.com/goda6565/nexus-user-auth/interface/handler/user/authentication"
	profileHandler "github.com/goda6565/nexus-user-auth/interface/handler/user/profile"
	registrationHandler "github.com/goda6565/nexus-user-auth/interface/handler/user/registration"
	"github.com/goda6565/nexus-user-auth/interface/middleware"
	"github.com/goda6565/nexus-user-auth/pkg/logger"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

type ServerInterfaceImpl struct {
	*registrationHandler.UserRegistrationHandler
	*authenticationHandler.UserAuthenticationHandler
	*profileHandler.UserProfileHandler
}

// swagger設定
func setUpSwagger(router *gin.Engine) (*openapi3.T, error) {
	swagger, err := gen.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := utils.GetEnvDefault("ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InfoInstanceName, SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	router := gin.New()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	swagger, err := setUpSwagger(router)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())
	router.GET("/health", handler.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(10 * time.Second))
		v1 := apiGroup.Group("/v1")

		v1.Use(middleware.AuthMiddleware("/api/v1/profile"))

		// OapiRequestValidator は v1 グループに適用（認証は後述の動的ミドルウェアで行う）
		v1.Use(ginMiddleware.OapiRequestValidatorWithOptions(swagger, &ginMiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
					// ここでは何もせず nil を返す
					return nil
				},
			},
		}))

		// すべてのハンドラーをひとつにまとめる
		userRepositoryImpl := repository.NewUserRepository(db)
		userRegistrationService := registrationService.NewUserRegistrationService(userRepositoryImpl)
		userRegistrationHandler := registrationHandler.NewUserRegistrationHandler(userRegistrationService)
		userAuthenticationService := authenticationService.NewUserAuthenticationService(userRepositoryImpl)
		userAuthenticationHandler := authenticationHandler.NewUserAuthenticationHandler(userAuthenticationService)
		userProfileService := profileService.NewUserProfileService(userRepositoryImpl)
		userProfileHandler := profileHandler.NewUserProfileHandler(userProfileService)

		serverInterface := &ServerInterfaceImpl{
			UserRegistrationHandler:   userRegistrationHandler,
			UserAuthenticationHandler: userAuthenticationHandler,
			UserProfileHandler:        userProfileHandler,
		}

		// v1 グループにハンドラーを登録する
		gen.RegisterHandlers(v1, serverInterface)
	}
	return router, nil
}
