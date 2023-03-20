package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tipbk/doodle/model"
)

type jwtService struct {
	jwtSecret string
}

func NewJWTService(jwtSecret string) JWTService {
	return &jwtService{
		jwtSecret: jwtSecret,
	}
}

type JWTService interface {
	GenerateUserToken(u *model.User) (string, *time.Time, error)
	ValidateToken(tokenString string) (string, error)
}

func (s *jwtService) GenerateUserToken(u *model.User) (string, *time.Time, error) {
	expTime := time.Now().Add(time.Hour * 24)
	claims := jwt.MapClaims{}
	claims["jti"] = u.ID
	claims["email"] = u.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = expTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(s.jwtSecret)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", nil, err
	}
	return tokenString, &expTime, nil
}

func (s *jwtService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Replace with your own secret key
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("cannot parse claims")
	}

	id, ok := claims["jti"].(string)
	if !ok {
		return "", errors.New("cannot get id from claims")
	}

	return id, nil
}
