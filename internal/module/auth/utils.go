package auth

import (
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateTokens(userID string, key paseto.V4SymmetricKey) (string, string) {
	now := time.Now()

	access := paseto.NewToken()
	access.SetSubject(userID)
	access.SetIssuedAt(now)
	access.SetNotBefore(now)
	access.SetExpiration(now.Add(accessTokenTTL))
	access.SetString("typ", "access")

	refresh := paseto.NewToken()
	refresh.SetSubject(userID)
	refresh.SetIssuedAt(now)
	refresh.SetNotBefore(now)
	refresh.SetExpiration(now.Add(refreshTokenTTL))
	refresh.SetString("typ", "refresh")

	return access.V4Encrypt(key, nil), refresh.V4Encrypt(key, nil)
}
