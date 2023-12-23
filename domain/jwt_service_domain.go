package domain

import "github.com/dgrijalva/jwt-go"

//go:generate mockgen -source=jwt_service_domain.go -destination=./mock/jwt_service_domain.go
type JWTService interface {
	GenerateNewToken(username, role string) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
	RefreshToken(currenToken *jwt.Token) (string, error)
}
