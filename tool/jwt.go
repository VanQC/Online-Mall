package tool

import (
	"cloudrestaurant/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserId int
	jwt.StandardClaims
}

// 定义jwt加密的秘钥
// 密钥应该妥善保管，并且不应该直接硬编码在代码中，而是从安全的配置文件或环境变量中读取
var jwtKey = []byte("a_secret_create") // 秘钥具体内容根据需求设置

// GenerateToken 根据用户信息生成token
func GenerateToken(user *model.Member) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 有效期，设置的七天

	// 创建载荷(payload)实例
	claims := &Claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000, // 签名生效时间
			ExpiresAt: expirationTime.Unix(),    // 签名过期时间
			IssuedAt:  time.Now().Unix(),        // 发布时间
			Issuer:    user.UserName,            // 签名发布者
			Subject:   "user token",
		},
	}
	// 根据加密算法以及载荷初步生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 将初步生成的token结合，jwt秘钥，得到最终tokenString
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claim := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	return token, claim, err
}
