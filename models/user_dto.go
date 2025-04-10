package models

/**
 * @File: user_dto.go.go
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
	UserName string `json:"username"`
	Role     string `json:"role"`
}
