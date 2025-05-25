package db

import (
	"context"
	"fmt"
	"sync"
	"video_douyin/dal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB        *gorm.DB
	Redis     *redis.Client
	redisOnce sync.Once
)

func InitMySQL(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用表名复数化
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	// 自动迁移
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	return nil
}

func InitRedis(addr string, db int) error {
	var initErr error
	redisOnce.Do(func() {
		Redis = redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   db,
		})
		if _, err := Redis.Ping(context.Background()).Result(); err != nil {
			initErr = fmt.Errorf("redis连接失败: %w", err)
		}
	})
	return initErr
}
