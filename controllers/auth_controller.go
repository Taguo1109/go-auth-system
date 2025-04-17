package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
// @Param user body models.UserRegisterDTO true "使用者資訊（Email、Username、Password、Role 為必填）"
// @Success 200 {object} utils.JsonResult
// @Failure 400 {object} utils.JsonResult
// @Failure 500 {object} utils.JsonResult
// @Router /register [post]
func Register(c *gin.Context) {
	var input models.UserRegisterDTO

	if err := c.ShouldBindJSON(&input); err != nil {

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errFields := utils.ExtractFieldErrorMessages(input, ve)
			utils.ReturnError(c, utils.CodeParamInvalid, errFields, "欄位驗證失敗")
			return
		}
		utils.ReturnError(c, utils.CodeParamInvalid, err.Error())
		return
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = string(hashedPassword)

	// 將 DTO 映射成 DB entity
	userEntity := models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	// 建立使用者
	result := config.DB.Create(&userEntity)
	if result.Error != nil {
		utils.ReturnError(c, utils.CodeEmailExists, "該用戶已存在")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": input.Email + " :User registered successfully!",
		"user": gin.H{
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
	var input models.UserLoginDTO
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
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		Username:     dbUser.Username,
		Role:         dbUser.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	var responseDTO models.UserLoginResponseDTO
	bytes, _ := json.Marshal(safeUser)
	err = json.Unmarshal(bytes, &responseDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.JsonResult{
			StatusCode: "500",
			Msg:        "Failed to unmarshal JSON",
			MsgDetail:  "JSON轉換失敗",
		})
		return
	}
	userBytes, _ := json.Marshal(responseDTO)
	config.RDB.Set(config.Ctx, cacheKey, userBytes, 10*time.Minute)

	utils.ReturnSuccess(c, safeUser, "Login successful")
}

// RefreshToken 重新獲取 Token
// @Summary 使用者重新獲得 Token
// @Description 傳入 refresh_token 取得新的 access_token 與 refresh_token
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.JsonResult
// @Failure 401 {object} utils.JsonResult
// @Router /refresh [post]
func RefreshToken(c *gin.Context) {
	// 從 JSON 或 localStorage 帶進來
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		utils.ReturnError(c, utils.CodeParamInvalid, nil, "請提供 refresh_token")
		return
	}

	// 解析 token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	// 檢查Redis 是否存在黑名單
	isBlacklisted, _ := config.RDB.Exists(config.Ctx, "blacklist:refresh_token:"+input.RefreshToken).Result()
	if isBlacklisted == 1 {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "refresh_token 已失效，請重新登入")
		return
	}

	// 驗證 token 是否有效 & 是 refresh token
	if err != nil || !token.Valid || claims["token_type"] != "refresh" {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "refresh_token 無效或過期")
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "Token 內容無效")
		return
	}

	// 查詢使用者資料
	var dbUser models.User
	result := config.DB.Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "找不到使用者")
		return
	}

	// 產生新 token
	newAccessToken, newRefreshToken, err := utils.GenerateJWT(dbUser.Email, dbUser.ID, "User")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.JsonResult{
			StatusCode: "500",
			Msg:        "Can't generate access token",
			MsgDetail:  "無法產生新 token",
		})
		return
	}

	// 統一回傳 DTO
	safeUser := models.UserDTO{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		Username:     dbUser.Username,
		Role:         dbUser.Role,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	utils.ReturnSuccess(c, safeUser, "Token refreshed successfully")
}

// LogoutHandler 登出
// @Summary 使用者登出
// @Description 清除使用者的 access_token 和 refresh_token cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.JsonResult "成功登出訊息"
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {

	var input models.UserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ReturnError(c, utils.CodeBadRequest, nil, "格式錯誤")
		return
	}
	// 1️⃣ access_token 放進黑名單
	accessClaims, err := utils.ParseToken(input.AccessToken)
	if err == nil {
		if exp, ok := accessClaims["exp"].(float64); ok {
			ttl := time.Until(time.Unix(int64(exp), 0))
			// 若還沒過期則加入黑名單
			if ttl > 0 {
				config.RDB.Set(c, "blacklist:access_token:"+input.AccessToken, "1", ttl)
			}
		}
	}
	// 2️⃣ refresh_token 放進黑名單
	refreshClaims, err := utils.ParseToken(input.RefreshToken)
	if err == nil {
		if exp, ok := refreshClaims["exp"].(float64); ok {
			ttl := time.Until(time.Unix(int64(exp), 0))
			// 若還沒過期則加入黑名單
			if ttl > 0 {
				config.RDB.Set(c, "blacklist:refresh_token:"+input.RefreshToken, "1", ttl)
			}
		}
	}
	utils.ReturnSuccess(c, nil, "Logout successful")
}
