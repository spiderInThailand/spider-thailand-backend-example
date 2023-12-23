package jwt_service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	SecretKey  string
	ExpireTime time.Duration
	Issure     string
}

type AuthCustomClaims struct {
	Username string
	Role     string
	jwt.StandardClaims
}

func NewJWTService(secretKey string, expireTime time.Duration, issure string) *JWTService {
	return &JWTService{
		SecretKey:  secretKey,
		ExpireTime: expireTime,
		Issure:     issure,
	}
}

func (service *JWTService) GenerateNewToken(username, role string) (string, error) {
	claims := &AuthCustomClaims{
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(service.ExpireTime).Unix(),
			Issuer:    service.Issure,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(service.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (service *JWTService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])

		}
		return []byte(service.SecretKey), nil
	})
}

func (service *JWTService) RefreshToken(currenToken *jwt.Token) (string, error) {
	claims := currenToken.Claims.(jwt.MapClaims)

	username := fmt.Sprint(claims["Username"])
	role := fmt.Sprint(claims["Role"])

	newClaims := &AuthCustomClaims{
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(service.ExpireTime).Unix(),
			Issuer:    service.Issure,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenStr, err := token.SignedString([]byte(service.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
