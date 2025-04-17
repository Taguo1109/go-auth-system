package routes

/**
 * @File: routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午2:39
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-auth-system/docs"
	"go-auth-system/middlewares"
	"go-auth-system/routes/modules"
	"time"
)

func SetupRouter(r *gin.Engine) {
	// CORS中間件
	setupCorsMiddleware(r)
	// 加入全域錯誤攔截器
	r.Use(middlewares.GlobalErrorHandler())
	// Swagger文檔路由
	setupSwaggerRoutes(r)
	// 白名單(public)
	modules.SetPublicRoutes(r)
	// 全域錯誤測試(err)
	modules.SetErrRoutes(r)
	// 需權限的名單
	//(user)
	modules.SetUserRoutes(r)

}

// setupSwaggerRoutes 設置Swagger文檔路由
func setupSwaggerRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// setupCorsMiddleware 設置CORS中間件
func setupCorsMiddleware(r *gin.Engine) {
	// ✅ 加入 CORS 設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",       // 本地開發用
			"https://taguo1109.github.io", // GitHub Pages 正式站
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 如果你有用 cookie/token
		MaxAge:           12 * time.Hour,
	}))
}
