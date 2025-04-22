package services

import (
	"github.com/gin-gonic/gin"
	"go-auth-system/config"
	"go-auth-system/models"
	"go-auth-system/utils"
	"golang.org/x/crypto/bcrypt"
)

/**
 * @File: auth_service.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/22 上午9:28
 * @Software: GoLand
 * @Version:  1.0
 */

// AuthService encapsulates authentication logic
type AuthService struct{}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register 直接在 Service 呼叫 utils 返回 JSON，無需回傳任何參數
func (s *AuthService) Register(c *gin.Context, input models.UserRegisterDTO) {
	// 1. 密碼雜湊
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ReturnError(c, utils.CodeServerError, nil, "密碼加密失敗")
		return
	}

	// 2. 準備用戶實體
	user := models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hashed),
		Role:     input.Role,
	}

	// 3. 寫入資料庫
	if err := config.DB.Create(&user).Error; err != nil {
		// 假設 config.ErrDuplicateKey 代表主鍵衝突
		utils.ReturnError(c, utils.CodeEmailExists, nil, "該用戶已存在")
		return
	}

	// 4. 組裝 DTO
	resp := models.UserLoginResponseDTO{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}

	// 5. 回傳成功 JSON
	utils.ReturnSuccess(c, resp, user.Email+" :註冊成功")
}
