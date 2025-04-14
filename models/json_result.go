package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @File: json_result.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/14 下午2:55
 * @Software: GoLand
 * @Version:  1.0
 */

type JsonResult struct {
	StatusCode int         `json:"status_code"`
	Msg        interface{} `json:"msg"`
	Data       interface{} `json:"data"`
}

type JsonErrorResult struct {
	StatusCode int         `json:"status_code"`
	Msg        interface{} `json:"msg"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}) {
	json := &JsonResult{StatusCode: code, Msg: msg, Data: data}
	c.JSON(http.StatusOK, json)
}
func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonErrorResult{StatusCode: code, Msg: msg}
	c.JSON(http.StatusOK, json)
}
