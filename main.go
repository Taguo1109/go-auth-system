package main

/**
 * @File: main.go
 * @Description:
 *
 *	專案入口8080
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午2:37
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go-auth-system/config"
	"go-auth-system/middlewares"
	"go-auth-system/models"
	"go-auth-system/routes"
	"log"
)

// @title           登入系統API
// @version         1.0
// @description     登入系統的RESTful API接口文檔
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 請求頭中必須添加 Authorization Bearer {token}，例如 "Authorization: Bearer abcxyz"
func main() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 註冊密碼驗證器
		_ = v.RegisterValidation("pwd_validation", middlewares.UserPwd)
		// 註冊使用者名稱驗證器
		_ = v.RegisterValidation("username_validation", middlewares.UserName)
	}
	// DB初始化
	config.ConnectDB()
	// Redis 初始化
	config.InitRedis()

	// 自動Create User Table
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := gin.Default()
	routes.SetupRouter(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
