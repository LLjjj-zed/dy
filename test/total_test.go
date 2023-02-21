package test

import (
	"douyin.core/Model"
	"douyin.core/router"
	"douyin.core/service"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestAll(t *testing.T) {
	Model.InitDB_test()
	go service.RunMessageServer()
	r := gin.Default()
	router.InitRouter(r)
	r.Run("0.0.0.0:8080")
}
