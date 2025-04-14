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
	"go-auth-system/config"
	"go-auth-system/models"
	"go-auth-system/routes"
	"log"
)

func main() {

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
