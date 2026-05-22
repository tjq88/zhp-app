package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// success 按项目统一格式返回成功响应。
func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "success",
		"data": data,
	})
}

// fail 按项目统一格式返回失败响应。
func fail(c *gin.Context, status int, code, msg string) {
	c.JSON(status, gin.H{
		"code": code,
		"msg":  msg,
	})
}
