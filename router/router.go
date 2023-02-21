package router

import (
	"douyin.core/controller"
	"douyin.core/handler/Comment"
	"douyin.core/handler/Like"
	"douyin.core/handler/User"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")
	apiRouter.POST("/user/register/", User.UserRegistHandler)
	r.Use(middleware.JWTMiddleWare())
	//inteeraction
	//apiRouter.POST("/favorite/action/", Like.LikeHandler)
	//apiRouter.GET("/favorite/list/", Like.GetLikeList)
	//apiRouter.POST("/comment/action/", Comment.CommentActionHandler)
	//apiRouter.GET("/comment/list/", Comment.GetCommentList)
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	//apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

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
