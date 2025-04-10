package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-system/config"
	"go-auth-system/models"
	"go-auth-system/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

/**
 * @File: auth_controller.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午3:13
 * @Software: GoLand
 * @Version:  1.0
 */

// Register 會員註冊
func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = string(hashedPassword)

	// 建立使用者
	result := config.DB.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": input.Email + " :User registered successfully!",
		"user": gin.H{
			"id":       input.ID,
			"email":    input.Email,
			"userName": input.Username,
		},
	})
}

// Login 會員登入
func Login(c *gin.Context) {
	var input models.User
	var dbUser models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 找出該 Email 的使用者
	result := config.DB.Where("email = ?", input.Email).First(&dbUser)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email does not exist"})
		return
	}

	// 檢查密碼
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	// 產生 JWT token
	accessToken, refreshToken, err := utils.GenerateJWT(dbUser.Email, dbUser.ID, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"token":        accessToken,
		"refreshToken": refreshToken,
	})
}

// RefreshToken 重新獲取Token
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// 1️⃣ 驗證 refresh token
	claims := jwt.MapClaims{}
	fmt.Println("初始化claims:", claims)
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})
	fmt.Println("賦值後claims:", claims)

	// 檢查是否是refreshToken
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not a refresh token"})
		return
	}

	fmt.Print(err)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// 2️⃣ 解析出 email
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
		return
	}

	// 3️⃣ 用 email 去資料庫撈使用者資訊
	var dbUser models.User
	result := config.DB.Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// 4️⃣ 產生新的 access + refresh token
	newAccessToken, newRefreshToken, err := utils.GenerateJWT(dbUser.Email, dbUser.ID, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Token refreshed successfully",
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}
