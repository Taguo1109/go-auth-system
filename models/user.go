package models

/**
 * @File: user.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午2:55
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"gorm.io/gorm"
)

// User 建立User Table
type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
