package modules

import (
	"github.com/gin-gonic/gin"
	"go-auth-system/controllers"
	"go-auth-system/middlewares"
)

/**
 * @File: user_routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/17 上午10:35
 * @Software: GoLand
 * @Version:  1.0
 */

func SetUserRoutes(r *gin.Engine) {
	// 需權限的名單
	user := r.Group("/user")
	user.Use(middlewares.JWTAuthMiddleware())
	{
		user.GET("/profile", controllers.GetProfile)
	}
}
