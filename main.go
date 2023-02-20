package main

import (
	"douyin.core/config"
	"douyin.core/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()
	r := gin.Default()
	initRouter(r)
	r.Run(fmt.Sprintf(":%d", config.Info.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
