package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-auth-system/config"
	"go-auth-system/models"
	"net/http"
	"time"
)

/**
 * @File: user_controller.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/10 上午10:26
 * @Software: GoLand
 * @Version:  1.0
 */

func GetProfile(c *gin.Context) {
	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No email in token"})
		return
	}
	email := emailVal.(string)

	// 1️⃣ 先從 Redis 查快取
	cacheKey := "user:" + email
	cached, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		// 如果有快取，直接回傳
		var cachedUser models.UserDTO
		// json.Unmarshal 將資料JSON格式化
		if err := json.Unmarshal([]byte(cached), &cachedUser); err == nil {
			c.JSON(http.StatusOK, gin.H{"user": cachedUser, "from": "cache"})
			return
		}
	}

	// 2️⃣ 沒快取，查 DB
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3️⃣ 查到後，存入 Redis 快取（設 10 分鐘過期）
	safeUser := models.UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		UserName: user.Username,
		Role:     user.Role,
	}
	userBytes, _ := json.Marshal(safeUser)
	config.RDB.Set(config.Ctx, cacheKey, userBytes, 10*time.Minute)
	c.JSON(http.StatusOK, gin.H{"user": user, "from": "db"})
}
