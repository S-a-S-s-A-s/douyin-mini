package main

import (
	"douyin-mini/db"
	"douyin-mini/service"

	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	InitRouter(r)
	db.Init()
	//创建测试数据
	db.CreateData()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
