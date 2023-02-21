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
	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
