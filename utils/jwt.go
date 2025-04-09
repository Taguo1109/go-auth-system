package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

/**
 * @File: jwt.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午3:43
 * @Software: GoLand
 * @Version:  1.0
 */

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT 生成Token
func GenerateJWT(email string) (string, error) {
	fmt.Println("🔐 JWT_SECRET in Login =", os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
