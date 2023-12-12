package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type UserAuthService interface {
	GenerateToken(userID int, role int) (string, error)
	ValidasiToken(token string) (*jwt.Token, error)
}

var SecretKey []byte

type jwtService struct {
}

func NewUserAuthService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) SetSecretKey(key string) {
	SecretKey = []byte(key)
}

func (s *jwtService) GenerateToken(userID int, role int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID
	claim["role"] = role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *jwtService) ValidasiToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}
