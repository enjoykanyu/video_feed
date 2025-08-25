package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"video_douyin/dal/model"
)

const (
	// 用户兴趣标签缓存键前缀
	UserInterestKeyPrefix = "user:interest:"
	// 用户兴趣标签过期时间
	UserInterestExpire = 24 * time.Hour
	// 热门标签排行榜键
	HotTagsKey = "hot:tags"
	// 标签相似度缓存键前缀
	TagSimilarityKeyPrefix = "tag:similarity:"
)

// UserInterestCache 用户兴趣标签缓存结构
type UserInterestCache struct {
	UserID    int64              `json:"user_id"`
	Tags      map[string]float64 `json:"tags"` // tag_name -> weight
	UpdatedAt time.Time          `json:"updated_at"`
	Version   int64              `json:"version"` // 版本号，用于缓存失效
}

// CacheUserInterests 缓存用户兴趣标签
func CacheUserInterests(ctx context.Context, userID int64, interests []model.UserInterest) error {
	key := fmt.Sprintf("%s%d", UserInterestKeyPrefix, userID)

	// 构建标签权重映射
	tags := make(map[string]float64)
	for _, interest := range interests {
		tags[interest.TagName] = interest.Weight
	}

	cache := &UserInterestCache{
		UserID:    userID,
		Tags:      tags,
		UpdatedAt: time.Now(),
		Version:   time.Now().Unix(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal user interests: %v", err)
	}

	return RedisClient.Set(ctx, key, data, UserInterestExpire).Err()
}

// GetUserInterests 获取用户兴趣标签缓存
func GetUserInterests(ctx context.Context, userID int64) (*UserInterestCache, error) {
	key := fmt.Sprintf("%s%d", UserInterestKeyPrefix, userID)

	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cache UserInterestCache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user interests: %v", err)
	}

	return &cache, nil
}

// UpdateUserInterestWeight 实时更新用户兴趣标签权重
func UpdateUserInterestWeight(ctx context.Context, userID int64, tagName string, weightDelta float64) error {
	key := fmt.Sprintf("%s%d", UserInterestKeyPrefix, userID)

	// 使用Redis事务确保原子性
	txf := func(tx *redis.Tx) error {
		// 获取当前缓存
		data, err := tx.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		var cache UserInterestCache
		if err == redis.Nil {
			// 缓存不存在，创建新的
			cache = UserInterestCache{
				UserID:    userID,
				Tags:      make(map[string]float64),
				UpdatedAt: time.Now(),
				Version:   time.Now().Unix(),
			}
		} else {
			if err := json.Unmarshal([]byte(data), &cache); err != nil {
				return fmt.Errorf("failed to unmarshal cache: %v", err)
			}
		}

		// 更新权重
		currentWeight := cache.Tags[tagName]
		newWeight := currentWeight + weightDelta
		if newWeight < 0 {
			newWeight = 0
		}
		cache.Tags[tagName] = newWeight
		cache.UpdatedAt = time.Now()
		cache.Version = time.Now().Unix()

		// 序列化并保存
		newData, err := json.Marshal(cache)
		if err != nil {
			return fmt.Errorf("failed to marshal updated cache: %v", err)
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, newData, UserInterestExpire)
			return nil
		})
		return err
	}

	// 重试机制
	for i := 0; i < 3; i++ {
		err := RedisClient.Watch(ctx, txf, key)
		if err == nil {
			break
		}
		if err == redis.TxFailedErr {
			continue
		}
		return err
	}

	// 同时更新热门标签排行榜
	RedisClient.ZIncrBy(ctx, HotTagsKey, weightDelta, tagName)

	return nil
}

// GetHotTags 获取热门标签排行榜
func GetHotTags(ctx context.Context, count int64) ([]string, error) {
	result, err := RedisClient.ZRevRange(ctx, HotTagsKey, 0, count-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hot tags: %v", err)
	}
	return result, nil
}

// BatchUpdateUserInterests 批量更新用户兴趣标签
func BatchUpdateUserInterests(ctx context.Context, updates []UserInterestUpdate) error {
	pipe := RedisClient.Pipeline()

	for _, update := range updates {
		key := fmt.Sprintf("%s%d", UserInterestKeyPrefix, update.UserID)

		// 使用HINCRBY更新标签权重
		pipe.HIncrBy(ctx, key, update.TagName, int64(update.WeightDelta*1000)) // 转换为整数避免精度问题
		pipe.Expire(ctx, key, UserInterestExpire)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// UserInterestUpdate 用户兴趣标签更新结构
type UserInterestUpdate struct {
	UserID      int64   `json:"user_id"`
	TagName     string  `json:"tag_name"`
	WeightDelta float64 `json:"weight_delta"`
}

// InvalidateUserInterestCache 使缓存失效
func InvalidateUserInterestCache(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%s%d", UserInterestKeyPrefix, userID)
	return RedisClient.Del(ctx, key).Err()
}

// GetUserInterestVersion 获取用户兴趣标签版本号
func GetUserInterestVersion(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf("%s%d", UserInterestKeyPrefix, userID)

	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	var cache UserInterestCache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return 0, fmt.Errorf("failed to unmarshal cache: %v", err)
	}

	return cache.Version, nil
}
