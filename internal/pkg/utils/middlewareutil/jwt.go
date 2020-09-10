package middlewareutil

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	// UserSignedKey 定义UserSignedKey
	UserSignedKey = "ut_blogger_"
)

// MySecret 定义JWT TOKEN的加密盐
var MySecret = []byte(UserSignedKey)

// MyClaims 自定义Claims结构
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateAccessToken  生成jwt token
func CreateAccessToken(username string, expired int64) (string, error) {
	claims := &MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expired,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(MySecret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
