package config

import (
	"log"
	"time"

	"exchangeapp/global"
	"exchangeapp/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to configure database: %v", err)
	}
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns) // 最大打开连接数
	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)                    // 连接最大生命周期

	global.DB = db

	if err := db.AutoMigrate(&models.User{}, &models.ExchangeRate{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
