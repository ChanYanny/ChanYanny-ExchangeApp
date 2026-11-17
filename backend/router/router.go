package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}
	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRates)
	api.Use(middlewares.AuthMiddleware()) // 下面的路由都需要认证
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
		api.GET("/articles/:id", controllers.GetArticlesByID)
		api.GET("/articles", controllers.GetArticles)
		api.POST("/articles", controllers.CreateArticle)
		api.POST("/articles/:id/likes", controllers.LikeArticle)
		api.GET("/articles/:id/likes", controllers.GetArticleLikes)
	}
	return r
}
