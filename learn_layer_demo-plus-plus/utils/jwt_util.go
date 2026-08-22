// 创建时间：2026/8/21 下午6:06
package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// ============学习demo临时密钥，后期简历项目迁移到配置文件============
const jwtSecret = "demo-secret-key-2026-go-gin-layer"
const jwtExpireHour = 2

// MyClaims JWT载荷
type MyClaims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 根据userId生成token字符串
func GenerateToken(userId uint) (string, error) {
	// 1.把传入的 userId 放进JWT载荷，同时设置签发时间、过期时间
	claims := MyClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// 2.以HS256算法创建JWT对象，把载荷数据打包进去
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 3.用后端密钥签名，生成最终token字符串返回给前端
	return tokenObj.SignedString([]byte(jwtSecret))
}

// ParseToken 解析token，校验签名、过期时间，返回userId
func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims.UserId, nil
	}
	return 0, fmt.Errorf("token无效")
}
