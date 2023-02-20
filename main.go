package main

import (

	"douyin.core/Model"
	"douyin.core/service"
	"github.com/gin-gonic/gin"
)

func main() {

	Model.InitDB()
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)
	r.Run(GinSocket) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
