package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-system/utils"
	"net/http"
)

/**
 * @File: jwt_middleware.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/9 上午11:17
 * @Software: GoLand
 * @Version:  1.0
 */

// JWTAuthMiddleware 是一個 Gin middleware，用來攔截需要登入驗證的 API
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1️⃣ 從 Header 中讀取 Authorization 欄位
		//authHeader := c.GetHeader("Authorization")
		// 1️⃣ 從 cookie 中讀取 access_token
		tokenString, err := c.Cookie("access_token")

		// 2️⃣ 檢查 Header 是否以 "Bearer " 開頭
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, utils.JsonResult{
				StatusCode: "401",
				Msg:        "No token provided in cookie",
				MsgDetail:  "請先登入或確認 Cookie 中是否正確設定 token",
			})
			c.Abort()
			return
		}

		// 4️⃣ 驗證並解析 JWT Token（ParseWithClaims 會同時做 signature 驗證與解析 payload）
		//    - claims：token 的內容（例如 email、userId、role）會被填入這個 map
		//    - jwt.ParseWithClaims：接收 token 字串、claims 接收變數、驗證密鑰 callback
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtKey, nil
		})

		// 5️⃣ 確認 token 有效，若驗證失敗則中止請求
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.JsonResult{StatusCode: "401", Msg: "Invalid token", MsgDetail: "Token 無效請確認"})
			c.Abort()
			return
		}

		// 6️⃣ 直接使用解析出的 claims（已是 jwt.MapClaims）
		//     - 將 email 設入 Gin Context 傳遞給後面 handler 使用
		c.Set("email", claims["email"])
		c.Set("userId", claims["userId"])
		c.Set("role", claims["role"])

		// 7️⃣ 這一行是放行（讓後面的 handler 繼續執行）
		c.Next()
	}
}
