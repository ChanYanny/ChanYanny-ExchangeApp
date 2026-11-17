package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、 DeletedAt（软删除）字段
	Title      string `binding:"required"`
	Content    string `binding:"required"`
	Preview    string `binding:"required"`
	Likes      int    `gorm:"default:0"`
}
