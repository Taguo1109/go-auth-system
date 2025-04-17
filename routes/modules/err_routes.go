package modules

import (
	"github.com/gin-gonic/gin"
	"go-auth-system/controllers"
)

/**
 * @File: err_routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/17 上午10:35
 * @Software: GoLand
 * @Version:  1.0
 */

func SetErrRoutes(r *gin.Engine) {

	// 全域錯誤測試
	errTest := r.Group("/err")
	{
		errTest.GET("assertion-panic", controllers.AssertionPanic)
		errTest.GET("slice-panic", controllers.SlicePanic)
		errTest.GET("nil-panic", controllers.NilPanic)
		errTest.GET("test-panic", controllers.TestPanic)
	}
}
