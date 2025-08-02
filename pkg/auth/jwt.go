package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserID uint64 `json:"user_id"`
	RoleID string `json:"role_id"`
}

type Claims struct {
	UserID  uint64 `json:"user_id"`
	RoleID  string `json:"role_id"`
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret")                // use config
var AccessTokenExpiry = time.Hour * 24       // 24 hours
var RefreshTokenExpiry = time.Hour * 24 * 30 // 30 days

func Parse(token *jwt.Token) (interface{}, error) {
	// Verify signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return 0, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	// Return the secret key used for signing
	return jwtKey, nil
}

// GenerateJWT generates a JWT token for a user.
func GenerateJWT(user *User) (string, error) {
	claims := &Claims{
		UserID: user.UserID,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// GenerateFreshToken generates a fresh token and stores it in Redis
func GenerateFreshToken(userID uint64) (string, error) {
	// Generate random tid
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tID := hex.EncodeToString(b)

	// Generate JWT token
	claims := &Claims{
		UserID:  userID,
		TokenID: tID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GetJWTKey() []byte {
	return jwtKey
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
