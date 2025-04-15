package models

/**
 * @File: user_dto.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/10 上午11:09
 * @Software: GoLand
 * @Version:  1.0
 */

type UserDTO struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterDTO struct {
	Email    string `gorm:"type:varchar(191);unique" json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}
