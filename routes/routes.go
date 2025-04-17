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
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-auth-system/docs"
	"go-auth-system/middlewares"
	"go-auth-system/routes/modules"
)

func SetupRouter(r *gin.Engine) {

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
