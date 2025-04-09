package controllers

import (
	"github.com/gin-gonic/gin"
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
	token, err := utils.GenerateJWT(dbUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token})
}
