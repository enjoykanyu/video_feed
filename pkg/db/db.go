package db

import (
	"context"
	"fmt"
	"video_douyin/dal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitMySQL(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	return nil
}

func InitRedis(addr, password string, db int) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
		panic(fmt.Sprintf("Redis连接失败: %v", err))
	}
}
