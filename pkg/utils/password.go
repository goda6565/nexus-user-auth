package utils

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/goda6565/nexus-user-auth/errs"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.NewPkgError("failed to hash password")
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
