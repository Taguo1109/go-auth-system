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
	"go-auth-system/controllers"
	"go-auth-system/middlewares"
)

func SetupRouter(r *gin.Engine) {

	// 白名單
	public := r.Group("/")
	{
		public.POST("/login", controllers.Login)
		public.POST("/register", controllers.Register)
		public.POST("/refresh", controllers.RefreshToken)
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	// 需權限的名單
	user := r.Group("/user")
	user.Use(middlewares.JWTAuthMiddleware())
	{
		user.GET("/profile", controllers.GetProfile)
	}

}
