package models

/**
 * @File: user.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午2:55
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"time"
)

// User 建立User Table
type User struct {
	ID        uint   `gorm:"primary"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
