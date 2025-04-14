package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

/**
 * @File: jwt.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午3:43
 * @Software: GoLand
 * @Version:  1.0
 */

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT 生成Token
func GenerateJWT(email string, userId uint, role string) (string, string, error) {
	fmt.Println("🔐 JWT_SECRET in Login =", os.Getenv("JWT_SECRET"))
	// 1️⃣ Access Token - 壽命短（2 小時）
	accessClaims := jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"role":   role,
		"exp":    time.Now().Add(2 * time.Hour).Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString(JwtKey)
	if err != nil {
		return "", "", err
	}

	// 2️⃣ Refresh Token - 壽命長（7 天）
	refreshClaims := jwt.MapClaims{
		"email":      email,
		"token_type": "refresh", // 來辨別refresh 提供Refresh的API使用
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString(JwtKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
