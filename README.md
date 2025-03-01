# Nexus User Auth

**Nexus User Auth** は、DDD（ドメイン駆動設計）に基づいて実装されたユーザー認証およびプロフィール管理サービスです。  
Go言語を用いて、ユーザー登録、ログイン、JWTを利用したトークンリフレッシュ、ユーザープロフィールの取得・更新・削除などの機能を提供します。

## 機能概要

- **ユーザー登録**  
  ユーザーの新規登録を行います。  
  - サービス: `UserRegistrationService`  
  - エンドポイント例: `POST /api/v1/register`

- **ユーザープロフィール管理**  
  ユーザー情報の取得、更新、削除を行います。  
  - サービス: `UserProfileService`  
  - エンドポイント例:
    - 取得: `GET /api/v1/profile`
    - 更新: `PUT /api/v1/profile`
    - 削除: `DELETE /api/v1/profile`  
    ※ これらのエンドポイントでは、JWTからユーザー情報（objIDなど）を取得するため、URLにユーザーIDを含める必要はありません。

- **ユーザー認証**  
  ユーザーログインとトークンリフレッシュの処理を提供します。  
  - サービス: `UserAuthenticationService`  
  - エンドポイント例:
    - ログイン: `POST /api/v1/auth/login`
    - トークンリフレッシュ: `POST /api/v1/auth/refresh`

- **JWT管理**  
  アクセストークンおよびリフレッシュトークンの生成・検証を行います。  
  - ユーティリティ: `pkg/utils`

## プロジェクト構成

```
.
├── Makefile
├── README.md
├── api
│   ├── config.yaml
│   └── openapi.yaml
├── application
│   └── service
│       └── user
│           ├── authentication
│           │   ├── user_authentication_service.go
│           │   └── user_authentication_service_test.go
│           ├── profile
│           │   ├── user_profile_service.go
│           │   └── user_profile_service_test.go
│           └── registration
│               ├── user_registration_service.go
│               └── user_registration_service_test.go
├── atlas.hcl
├── docker-compose.yaml
├── domain
│   ├── timeobj
│   │   ├── time_obj.go
│   │   └── time_obj_test.go
│   └── user
│       ├── entity
│       │   ├── user_entity.go
│       │   └── user_entity_test.go
│       ├── repository
│       │   └── user_repository.go
│       └── value
│           ├── user_avatar_url.go
│           ├── user_avatar_url_test.go
│           ├── user_email.go
│           ├── user_email_test.go
│           ├── user_obj_id.go
│           ├── user_obj_id_test.go
│           ├── user_password.go
│           ├── user_password_test.go
│           ├── user_role.go
│           ├── user_role_test.go
│           ├── user_username.go
│           └── user_username_test.go
├── errs
│   ├── domain.go
│   ├── infra.go
│   ├── interface.go
│   ├── pkg.go
│   └── service.go
├── go.mod
├── go.sum
├── infrastructure
│   ├── database
│   │   ├── adapter
│   │   │   └── user_adapter.go
│   │   ├── config.go
│   │   ├── factory.go
│   │   ├── models
│   │   │   └── user_model.go
│   │   └── repository
│   │       ├── user_repository_impl.go
│   │       └── user_repository_impl_test.go
│   └── web
│       ├── config.go
│       ├── factory.go
│       └── gin.go
├── interface
│   ├── gen
│   │   └── api.go
│   ├── handler
│   │   ├── health.go
│   │   ├── health_test.go
│   │   └── user
│   │       ├── authentication
│   │       │   ├── user_authentication_handler.go
│   │       │   └── user_authentication_handler_test.go
│   │       ├── profile
│   │       │   ├── user_profile_handler.go
│   │       │   └── user_profile_handler_test.go
│   │       └── registration
│   │           ├── user_registration_handler.go
│   │           └── user_registration_handler_test.go
│   ├── keys
│   │   └── keys.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── timeout.go
│   └── router
│       └── router.go
├── main.go
├── migrations
│   ├── 20250301140523.sql
│   └── atlas.sum
└── pkg
    ├── logger
    │   └── logger.go
    ├── tester
    │   └── sqlite_suite.go
    └── utils
        ├── env.go
        ├── env_test.go
        ├── jwt.go
        ├── jwt_test.go
        ├── password.go
        └── password_test.go
```