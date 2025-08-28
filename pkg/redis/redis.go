package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	RedisClient *redis.Client
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 测试连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

// CacheVideoInfo 缓存视频信息
func CacheVideoInfo(ctx context.Context, videoID int64, data string) error {
	return RedisClient.Set(ctx, getVideoKey(videoID), data, 0).Err()
}

// GetVideoInfo 获取缓存的视频信息
func GetVideoInfo(ctx context.Context, videoID int64) (string, error) {
	return RedisClient.Get(ctx, getVideoKey(videoID)).Result()
}

// CacheUserInteractions 缓存用户互动信息
func CacheUserInteractions(ctx context.Context, userID int64, data string) error {
	return RedisClient.Set(ctx, getUserInteractionsKey(userID), data, 0).Err()
}

// GetUserInteractions 获取用户互动信息
func GetUserInteractions(ctx context.Context, userID int64) (string, error) {
	return RedisClient.Get(ctx, getUserInteractionsKey(userID)).Result()
}

// CacheHotVideos 缓存热门视频列表
func CacheHotVideos(ctx context.Context, data string) error {
	return RedisClient.Set(ctx, "hot:videos", data, 0).Err()
}

// GetHotVideos 获取热门视频列表
func GetHotVideos(ctx context.Context) (string, error) {
	return RedisClient.Get(ctx, "hot:videos").Result()
}

// CacheFollowing 缓存关注列表
func CacheFollowing(ctx context.Context, userID int64, data string) error {
	return RedisClient.Set(ctx, getFollowingKey(userID), data, 0).Err()
}

// GetFollowing 获取关注列表
func GetFollowing(ctx context.Context, userID int64) (string, error) {
	return RedisClient.Get(ctx, getFollowingKey(userID)).Result()
}

// 辅助函数：生成Redis键
func getVideoKey(videoID int64) string {
	return "video:" + string(videoID)
}

func getUserInteractionsKey(userID int64) string {
	return "user:" + string(userID) + ":interactions"
}

func getFollowingKey(userID int64) string {
	return "user:" + string(userID) + ":following"
}
