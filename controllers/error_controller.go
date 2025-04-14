package controllers

import "github.com/gin-gonic/gin"

/**
 * @File: error_controller.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/14 下午6:07
 * @Software: GoLand
 * @Version:  1.0
 */

// TestPanic 測試觸發 panic，用於驗證 GlobalErrorHandler。
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
	// panic("I am panic 😈")

	c.JSON(200, gin.H{"message": "看不到這行，前面就報錯了"})
}
