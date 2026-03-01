package util

import (
	"errors"
	"time"
	"go-vue-admin/global"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义声明结构体并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 我们这里需要额外记录一个username和userID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	RoleID   uint   `json:"roleId"`
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token已过期")
	TokenNotValidYet = errors.New("token尚未生效")
	TokenMalformed   = errors.New("非法token")
	TokenInvalid     = errors.New("无效token")
)

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.Config.JWT.SigningKey),
	}
}

// CreateToken 创建一个新的token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateClaims 创建Claims
func (j *JWT) CreateClaims(baseClaims CustomClaims) CustomClaims {
	expiresTime := time.Duration(global.Config.JWT.ExpiresTime) * time.Hour
	claims := CustomClaims{
		UserID:   baseClaims.UserID,
		Username: baseClaims.Username,
		RoleID:   baseClaims.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresTime)),
			Issuer:    global.Config.JWT.Issuer,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	return claims
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		// 处理不同类型的错误
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}
	return nil, TokenInvalid
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.JWT.ExpiresTime) * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
