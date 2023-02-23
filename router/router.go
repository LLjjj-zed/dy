package router

import (
	"douyin.core/controller"
	"douyin.core/handler/Comment"
	"douyin.core/handler/Like"
	"douyin.core/handler/User"
	"douyin.core/handler/Video"
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	//userRouter with no JWTMiddleWare
	userRouter := r.Group("/douyin/user")
	userRouter.POST("/register/", User.UserRegistHandler)
	userRouter.POST("/login/", User.UserLoginHandler)
	userRouter.GET("/", User.UserInfoHandler)

	//apiRouter use JWTMiddleWare
	apiRouter := r.Group("/douyin")
	apiRouter.Use(middleware.JWTMiddleWare())
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	//apiRouter.POST("/user/register/", controller.Register)

	apiRouter.POST("/publish/action/", Video.PublishVedioHandler)
	apiRouter.GET("/publish/list/", Video.UserPublishListHandler)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JWTMiddleWare(), Like.LikeHandler)
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
