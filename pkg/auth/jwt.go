package authx

import (
	"fmt"
	"time"

	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/golang-jwt/jwt/v4"
)

type JWTConf struct {
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
}

// GenJwtToken 生成 JWT 令牌
func GenJwtToken(secret string, expired int64, data map[string]any) (string, error) {
	// 创建声明
	claims := jwt.MapClaims{
		consts.Expire: time.Now().Add(time.Duration(expired) * time.Second).Unix(),
	}
	for key, value := range data {
		claims[key] = value
	}

	// 创建令牌
	var signs []*jwt.SigningMethodHMAC
	signs = append(signs, jwt.SigningMethodHS256)
	signs = append(signs, jwt.SigningMethodHS384)
	signs = append(signs, jwt.SigningMethodHS512)

	// 随机选择一种签名方法
	sign := signs[time.Now().Unix()%int64(len(signs))]

	token := jwt.NewWithClaims(sign, claims)

	// 签名生成token
	return token.SignedString([]byte(secret))
}

// ParseJwtToken 解析和验证 JWT 令牌
func ParseJwtToken(tokenStr, secret string) (jwt.MapClaims, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证令牌是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
