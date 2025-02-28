package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/goda6565/nexus-user-auth/errs"
)

// timeNowFunc は現在時刻取得用の関数。テスト用に差し替え可能にする。
var timeNowFunc = time.Now

type MyJWTClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("JWT_SECRET_KEY is not set")
	}
	return []byte(secret)
}

func GenerateTokens(objID string) (accessToken string, refreshToken string, err error) {
	// アクセストークン（短期有効）
	accessClaims := MyJWTClaims{
		ID: objID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ptf-auth-service",
			Subject:   objID,
			ExpiresAt: jwt.NewNumericDate(timeNowFunc().Add(24 * time.Hour)), // 24時間
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(getJWTSecret())
	if err != nil {
		return "", "", err
	}

	// リフレッシュトークン（長期有効）
	refreshClaims := MyJWTClaims{
		ID: objID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ptf-auth-service",
			Subject:   objID,
			ExpiresAt: jwt.NewNumericDate(timeNowFunc().Add(7 * 24 * time.Hour)), // 7日間
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(getJWTSecret())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

type TokenClaims struct {
	ObjID string
}

func ValidateToken(signedToken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// HMAC 系の署名アルゴリズムのみ許可
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.NewPkgError(fmt.Sprintf("jwt: unexpected signing method: %v", token.Header["alg"]))
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errs.NewPkgError("token signature is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.NewPkgError("token is expired")
		}
		return nil, errs.NewPkgError(fmt.Sprintf("jwt: %v", err))
	}

	claims, ok := token.Claims.(*MyJWTClaims)
	if !ok || !token.Valid {
		return nil, errs.NewPkgError("token is invalid")
	}

	return &TokenClaims{
		ObjID: claims.ID,
	}, nil
}

func ValidateRefreshToken(signedToken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.NewPkgError(fmt.Sprintf("jwt: unexpected signing method: %v", token.Header["alg"]))
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errs.NewPkgError("refresh token signature is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.NewPkgError("refresh token is expired")
		}
		return nil, errs.NewPkgError(fmt.Sprintf("jwt: %v", err))
	}

	claims, ok := token.Claims.(*MyJWTClaims)
	if !ok || !token.Valid {
		return nil, errs.NewPkgError("refresh token is invalid")
	}

	return &TokenClaims{
		ObjID: claims.ID,
	}, nil
}

func RefreshAccessToken(refreshToken string) (newAccessToken string, err error) {
	// リフレッシュトークンの検証
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 新しいアクセストークンを生成
	accessClaims := MyJWTClaims{
		ID: claims.ObjID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ptf-auth-service",
			Subject:   claims.ObjID,
			ExpiresAt: jwt.NewNumericDate(timeNowFunc().Add(24 * time.Hour)), // 24時間
		},
	}
	newAccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
