package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(viper.GetString("jwt.key"))

func GenerateToken(id uint) (string, error) {
	expireTime := time.Now().Add(time.Hour * 24) // 24小时token有效期

	//配置
	claims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("无效或过期的token")
	}
	return &claims, nil
}
