package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

type JwtService struct {
	secretKey string
	issuer    string
}

type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func NewJwtService() JwtService {
	return JwtService{
		secretKey: getSecretKey(),
		issuer:    "todo-service",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "sdfsdfsdfsd"
	}
	return secretKey
}

func (jwtService JwtService) GenerateToken(UserID string) string {
	log.Println(UserID)
	claims := &JwtCustomClaims{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    jwtService.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtService.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtService JwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected singing method %v", t.Header["alg"])
		}
		return []byte(jwtService.secretKey), nil
	})
}

func (jwtService JwtService) GetUserIdFromToken(token *jwt.Token) (string, error) {
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), nil
}
