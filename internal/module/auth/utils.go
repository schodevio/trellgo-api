package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func hashPassword(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func parseRefreshToken(refreshToken string, key paseto.V4SymmetricKey) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	token, err := parser.ParseV4Local(key, refreshToken, nil)
	if err != nil {
		return "", err
	}

	typ, err := token.GetString("typ")
	if err != nil || typ != "refresh" {
		return "", fmt.Errorf("invalid token type")
	}

	return token.GetSubject()
}
