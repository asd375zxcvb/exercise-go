package main

import (
	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 日志记录
	r.Use(gin.Logger())

	// 恢复
	r.Use(gin.Recovery())

	// CORS
	r.Use(cors.Default())

	r.Use(limits.RequestSizeLimiter(10))

	// 静态文件服务
	r.Static("/assets", "./assets")

	// 路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 运行服务器
	r.Run()
}
