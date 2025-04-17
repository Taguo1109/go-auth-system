package modules

import (
	"github.com/gin-gonic/gin"
	"go-auth-system/controllers"
)

/**
 * @File: public_routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/17 上午10:31
 * @Software: GoLand
 * @Version:  1.0
 */

func SetPublicRoutes(r *gin.Engine) {

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
}
