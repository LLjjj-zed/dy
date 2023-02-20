package main

import (
	"douyin.core/handler/Comment"
	"douyin.core/handler/Like"
	"douyin.core/handler/User"
	"douyin.core/handler/Video"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", Video.VideoFeedHandler)
	apiRouter.GET("/user/", User.UserInfoHandler)
	apiRouter.POST("/user/register/", User.UserRegistHandler)
	apiRouter.POST("/user/login/", User.UserLoginHandler)
	apiRouter.POST("/publish/action/", Video.PublishVedioHandler)
	apiRouter.GET("/publish/list/", Video.UserPublishListHandler)

	// extra apis - I
	apiRouter.POST("/favorite/action/", Like.LikeHandler)
	apiRouter.GET("/favorite/list/", Like.GetLikeList)
	apiRouter.POST("/comment/action/", Comment.CommentActionHandler)
	apiRouter.GET("/comment/list/", Comment.GetCommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
