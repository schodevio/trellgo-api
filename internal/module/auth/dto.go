package auth

import "time"

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
