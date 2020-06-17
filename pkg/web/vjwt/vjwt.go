package vjwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// 一些常量
var (
	// TokenExpired token过期
	TokenExpired error = errors.New("Token is expired")
	// TokenNotValidYet 令牌失效
	TokenNotValidYet error = errors.New("Token not active yet")
	// TokenMalformed 令牌格式错误
	TokenMalformed error = errors.New("That's not even a token")
	// TokenInvalid 令牌无效
	TokenInvalid error = errors.New("Couldn't handle this token:")
)

// New 新建一个jwt实例
func New() *JWT {
	return &JWT{}
}

// JWT 签名结构
type JWT struct {
	SignKey []byte
}

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID uint64 `json:"user_id"`
	jwt.StandardClaims
}

// GetSignKey 获取signKey
func (j *JWT) GetSignKey() string {
	return string(j.SignKey)
}

// SetSignKey 这是SignKey
func (j *JWT) SetSignKey(key string) *JWT {
	j.SignKey = []byte(key)
	return j
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
