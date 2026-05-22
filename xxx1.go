package zhp_app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Name string
	Age  int
}

func main() {
	//创建路由
	r := gin.Default()
	//绑定路由规则
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
	})
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "name": name, "action": action})
	})
	r.GET("/test", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "name": name, "age": age})
	})
	r.POST("/game", func(c *gin.Context) {
		var user User
		c.ShouldBindBodyWithJSON(&user)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "user": user})
	})
	//监听端口
	r.Run(":8080")
}
