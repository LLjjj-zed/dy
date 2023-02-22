package main

import (
	"douyin.core/Model"
	"douyin.core/handler/User"
	"douyin.core/handler/Video"
	"douyin.core/middleware"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
)

const GinSocket string = "0.0.0.0:8080"

func initRouter(r *gin.Engine) {
	Model.InitDB_test()
	//Model.InitDB()
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", Video.VideoFeedHandler)
	apiRouter.GET("/user/", middleware.JWTMiddleWare(), User.UserInfoHandler)
	apiRouter.POST("/user/register/", User.UserRegistHandler)
	apiRouter.POST("/user/login/", User.UserLoginHandler)
	apiRouter.POST("/publish/action/", Video.PublishVedioHandler)
	apiRouter.GET("/publish/list/", Video.UserPublishListHandler)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - IIz`
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
