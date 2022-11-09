package auth

import "github.com/golang-jwt/jwt"

type GetClaimsRequest struct {
	Token string `json:"access_token"`
}
type GetClaimsResponse struct {
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
}

type Jwt struct {
	AccessToken  string
	AccessExp    int64
	RefreshToken string
	RefreshExp   int64
}

type Claims struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
}

type ValidateRequest struct {
	Token string
}

type RenewJwtRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RenewJwtResponse struct {
	AccessToken  string `json:"access_token"`
	AccessExp    int64  `json:"access_exp"`
	RefreshToken string `json:"refresh_token"`
	RefreshExp   int64  `json:"refresh_exp"`
}
