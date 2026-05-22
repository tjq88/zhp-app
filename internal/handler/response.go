package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "success",
		"data": data,
	})
}

func fail(c *gin.Context, status int, code, msg string) {
	c.JSON(status, gin.H{
		"code": code,
		"msg":  msg,
	})
}
