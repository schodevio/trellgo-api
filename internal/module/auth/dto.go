package auth

import "time"

// Errors
const (
	INVALID_CREDENTIALS         string = "invalid credentials"
	INVALID_OR_EXPIRED_TOKEN    string = "invalid or expired token"
	REFRESH_TOKEN_NOT_FOUND     string = "refresh token not found"
	REFRESH_TOKEN_REVOKE_FAILED string = "failed to revoke refresh token"
	PASSWORD_HASH_FAILED        string = "failed to hash password"
	SIGN_IN_FAILED              string = "failed to sign in"
	SIGN_OUT_FAILED             string = "failed to sign out"
	SIGN_UP_FAILED              string = "failed to sign up"
	TOKEN_GENERATE_FAILED       string = "failed to generate tokens"
	USER_ALREADY_EXISTS         string = "user already exists"
	USER_CREATE_FAILED          string = "failed to create user"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SignInRequest struct {
	Email     string `json:"email"    validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string `json:"-"`
	IpAddress string `json:"-"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
}

type RefreshRequest struct {
	Token     string `json:"-"`
	UserAgent string `json:"-"`
	IpAddress string `json:"-"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignUpResponse struct {
	Email string `json:"email"`
}

// private

type generateTokensData struct {
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	IpAddress string `json:"ip_address"`
}
