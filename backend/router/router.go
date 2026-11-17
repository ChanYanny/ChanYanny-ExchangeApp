package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
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
