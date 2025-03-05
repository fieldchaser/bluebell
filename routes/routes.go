package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"web_framework/controllers"
	"web_framework/logger"
	"web_framework/middlewares"
)

func Setup() *gin.Engine {
	if viper.GetString("gin.mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//前端发来的http请求对应的路由操作
	//1.注册界面
	v1.POST("/signup", controllers.SignUp)
	v1.POST("/login", controllers.Login)

	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandlers)
		v1.GET("/community/:id", controllers.CommunityDetailHandlers)

		v1.POST("/post", controllers.CreatePostHandlers)
		v1.GET("/post/:id", controllers.GetPostDetailHandlers)
		v1.GET("/posts", controllers.GetPostListHandlers)

		v1.POST("/vote", controllers.PostVoteHandlers)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}
