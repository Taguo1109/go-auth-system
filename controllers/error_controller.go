package controllers

import (
	"github.com/gin-gonic/gin"
)

/**
 * @File: error_controller.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/14 下午6:07
 * @Software: GoLand
 * @Version:  1.0
 */

// AssertionPanic 測試1
// @Summary 測試 panic
// @Description 故意觸發 panic 來驗證全域錯誤攔截器是否生效
// @Tags Debug
// @Produce json
// @Success 200 {object} utils.JsonResult
// @Failure 500 {object} utils.JsonResult
// @Router /err/assertion-panic [get]
func AssertionPanic(c *gin.Context) {
	// ✅ 測試 #1：interface type assertion panic
	var x any = 123
	_ = x.(string) // runtime panic: interface conversion: int is not string

	c.JSON(200, gin.H{"message": "看不到這行，前面就報錯了"})
}

// SlicePanic 測試2
// @Summary 測試 panic
// @Description 故意觸發 panic 來驗證全域錯誤攔截器是否生效
// @Tags Debug
// @Produce json
// @Success 200 {object} utils.JsonResult
// @Failure 500 {object} utils.JsonResult
// @Router /err/slice-panic [get]
func SlicePanic(c *gin.Context) {

	// ✅ 測試 #2：slice 越界
	arr := []int{1, 2}
	_ = arr[10] // panic: index out of range

	c.JSON(200, gin.H{"message": "看不到這行，前面就報錯了"})

}

// NilPanic 測試3
// @Summary 測試 panic
// @Description 故意觸發 panic 來驗證全域錯誤攔截器是否生效
// @Tags Debug
// @Produce json
// @Success 200 {object} utils.JsonResult
// @Failure 500 {object} utils.JsonResult
// @Router /err/nil-panic [get]
func NilPanic(c *gin.Context) {

	// ✅ 測試 #3：nil map 寫入
	var m map[string]string
	m["key"] = "value" // panic: assignment to entry in nil map

	c.JSON(200, gin.H{"message": "看不到這行，前面就報錯了"})
}

// TestPanic 測試觸發 panic，用於驗證 GlobalErrorHandler。
// @Summary 測試 panic
// @Description 故意觸發 panic 來驗證全域錯誤攔截器是否生效
// @Tags Debug
// @Produce json
// @Success 200 {object} utils.JsonResult
// @Failure 500 {object} utils.JsonResult
// @Router /err/test-panic [get]
func TestPanic(c *gin.Context) {
	// ✅ 測試 #1：interface type assertion panic
	//var x any = 123
	//_ = x.(string) // runtime panic: interface conversion: int is not string

	// ✅ 測試 #2：slice 越界
	arr := []int{1, 2}
	_ = arr[10] // panic: index out of range

	// ✅ 測試 #3：nil map 寫入
	// var m map[string]string
	// m["key"] = "value" // panic: assignment to entry in nil map

	// ✅ 測試 #4：非法 JSON 結構
	// type BadStruct struct {
	//     F func()
	// }
	// b := BadStruct{}
	// json.Marshal(b) // panic: unsupported type: func()

	// ✅ 測試 #5：手動 panic
	panic("I am panic 😈")

}
