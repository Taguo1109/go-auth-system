package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-system/config"
	"go-auth-system/models"
	"go-auth-system/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

/**
 * @File: auth_controller.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午3:13
 * @Software: GoLand
 * @Version:  1.0
 */

// Register 會員註冊
// @Summary 使用者註冊
// @Description 新增一個使用者帳號
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "使用者資訊"
// @Success 200 {object} utils.JsonResult
// @Failure 500 {object} utils.JsonResult
// @Router /register [post]
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
		utils.ReturnError(c, utils.CodeEmailExists, "該用戶已存在")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": input.Email + " :User registered successfully!",
		"user": gin.H{
			"id":       input.ID,
			"email":    input.Email,
			"userName": input.Username,
			"role":     input.Role,
		},
	})
}

// Login 會員登入
// @Summary 使用者登入
// @Description 登入並取得 JWT Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.UserLoginDTO true "登入資訊"
// @Success 200 {object} utils.JsonResult
// @Failure 401 {object} utils.JsonResult
// @Router /login [post]
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

	// 快取使用者資料
	cacheKey := "user:" + dbUser.Email

	// 使dbUser 變成 JSON 標準格式 ，safeUser 存取需要的資訊進去
	safeUser := models.UserDTO{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Username: dbUser.Username,
		Role:     dbUser.Role,
	}
	userBytes, _ := json.Marshal(safeUser)
	config.RDB.Set(config.Ctx, cacheKey, userBytes, 10*time.Minute)

	// 設置 SameSite 模式為 Lax
	c.SetSameSite(http.SameSiteLaxMode)
	// 設定 access token Cookie，過期時間以秒計算（2 小時）
	c.SetCookie("access_token", accessToken, 2*60*60, "/", "localhost", true, true)
	// 設定 refresh token Cookie，過期 7 天
	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "localhost", true, true)

	utils.ReturnSuccess(c, nil, "Login successful")
}

// RefreshToken 重新獲取Token
// @Summary 重新獲取Token
// @Description 重新取得Token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.JsonResult
// @Failure 401 {object} utils.JsonResult
// @Router /register [post]
func RefreshToken(c *gin.Context) {

	// 1️⃣ 從 Cookie 中讀取 refresh token 驗證 refresh token
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, utils.JsonResult{
			StatusCode: "401",
			Msg:        "Missing refresh token in cookie",
			MsgDetail:  "找不到 Cookie 裡的refresh_token，請確認",
		})
		return
	}

	claims := jwt.MapClaims{}
	fmt.Println("初始化claims:", claims)
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
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

	// 設置 SameSite 模式為 Lax
	c.SetSameSite(http.SameSiteLaxMode)
	// 設置新的 Cookie，更新掉之前的 tokens
	c.SetCookie("access_token", newAccessToken, 2*60*60, "/", "localhost", true, true)
	c.SetCookie("refresh_token", newRefreshToken, 7*24*60*60, "/", "localhost", true, true)

	utils.ReturnSuccess(c, nil, "Token refreshed successfully")
}

// LogoutHandler 登出
// @Summary 使用者登出
// @Description 清除使用者的 access_token 和 refresh_token cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.JsonResult "成功登出訊息"
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {
	// 清除 access_token cookie
	c.SetCookie("access_token", "", -1, "/", "localhost", true, true)
	// 清除 refresh_token cookie
	c.SetCookie("refresh_token", "", -1, "/", "localhost", true, true)

	utils.ReturnSuccess(c, nil, "Logout successful")
}
