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
	"go-auth-system/controllers"
	_ "go-auth-system/docs"
	"go-auth-system/middlewares"
)

func SetupRouter(r *gin.Engine) {

	// 加入全域錯誤攔截器
	r.Use(middlewares.GlobalErrorHandler())
	// Swagger文檔路由
	setupSwaggerRoutes(r)

	// 白名單
	public := r.Group("/")
	{
		public.POST("/login", controllers.Login)
		public.POST("/register", controllers.Register)
		public.POST("/refresh", controllers.RefreshToken)
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		public.POST("/logout", controllers.LogoutHandler)
	}

	// 全域錯誤測試
	errTest := r.Group("/err")
	{
		errTest.GET("assertion-panic", controllers.AssertionPanic)
		errTest.GET("slice-panic", controllers.SlicePanic)
		errTest.GET("nil-panic", controllers.NilPanic)
		errTest.GET("test-panic", controllers.TestPanic)
	}

	// 需權限的名單
	user := r.Group("/user")
	user.Use(middlewares.JWTAuthMiddleware())
	{
		user.GET("/profile", controllers.GetProfile)
		// 測試 panic 路由
		user.GET("/test-panic", controllers.TestPanic)
	}

}

// setupSwaggerRoutes 設置Swagger文檔路由
func setupSwaggerRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
