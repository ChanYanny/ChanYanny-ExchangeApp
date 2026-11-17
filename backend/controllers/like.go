package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context) {
	// 实现点赞逻辑
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"
	if err := global.RedisClient.Incr(likeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article liked successfully"})
}

func GetArticleLikes(ctx *gin.Context) {
	// 实现获取文章点赞数逻辑
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"
	likes, err := global.RedisClient.Get(likeKey).Result()
	if err == redis.Nil {
		likes = "0"
	}else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}
