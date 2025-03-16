package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
	"web_framework/controllers"
	_ "web_framework/docs" // 千万不要忘了导入把你上一步生成的docs
	"web_framework/logger"
	"web_framework/middlewares"

	"github.com/swaggo/files"          // swagger embed files
	gs "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func Setup() *gin.Engine {
	if viper.GetString("gin.mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// CORS 配置，允许前端请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

		v1.GET("/posts2", controllers.GetPostListHandlers2)

		v1.POST("/vote", controllers.PostVoteHandlers)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	return r
}
