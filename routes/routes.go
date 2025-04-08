package routes

/**
 * @File: routes.go.go
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
)

func SetupRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/register", controllers.Register)
}
