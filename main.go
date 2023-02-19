package main

import (
	"douyin.core/service"
	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run("192.168.45.174:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
