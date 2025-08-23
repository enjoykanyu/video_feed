package behavior_processor

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
	"video_douyin/dal/model"
	"video_douyin/pkg/db"
	"video_douyin/pkg/rabbitmq"
	"video_douyin/pkg/redis"
)

// BehaviorProcessor 用户行为处理器
type BehaviorProcessor struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewBehaviorProcessor 创建行为处理器
func NewBehaviorProcessor() *BehaviorProcessor {
	return &BehaviorProcessor{}
}

// StartProcessing 开始处理用户行为事件
func (bp *BehaviorProcessor) StartProcessing() error {
	return rabbitmq.ConsumeUserBehavior(bp.handleUserBehavior)
}

// handleUserBehavior 处理用户行为事件
func (bp *BehaviorProcessor) handleUserBehavior(event *rabbitmq.UserBehaviorEvent) error {
	ctx := context.Background()

	klog.Infof("Processing user behavior: user_id=%d, video_id=%d, type=%s",
		event.UserID, event.VideoID, event.BehaviorType)

	// 1. 获取视频标签
	videoTags, err := bp.getVideoTags(ctx, event.VideoID)
	if err != nil {
		return fmt.Errorf("failed to get video tags: %v", err)
	}

	// 2. 计算权重增量
	weightIncrement := bp.calculateWeightIncrement(event)

	// 3. 实时更新Redis缓存
	if err := bp.updateRedisInterests(ctx, event.UserID, videoTags, weightIncrement); err != nil {
		klog.Errorf("Failed to update Redis interests: %v", err)
	}

	// 4. 异步更新MySQL数据库
	go func() {
		if err := bp.updateMySQLInterests(context.Background(), event); err != nil {
			klog.Errorf("Failed to update MySQL interests: %v", err)
		}
	}()

	// 5. 更新用户画像
	go func() {
		if err := bp.updateUserProfile(context.Background(), event); err != nil {
			klog.Errorf("Failed to update user profile: %v", err)
		}
	}()

	return nil
}

// getVideoTags 获取视频标签
func (bp *BehaviorProcessor) getVideoTags(ctx context.Context, videoID int64) ([]model.VideoTag, error) {
	var tags []model.VideoTag
	err := db.DB.Where("video_id = ?", videoID).Find(&tags).Error
	return tags, err
}

// calculateWeightIncrement 计算权重增量
func (bp *BehaviorProcessor) calculateWeightIncrement(event *rabbitmq.UserBehaviorEvent) float64 {
	var weight float64 = 0.1 // 基础权重

	// 根据观看时长调整权重
	if event.Duration > 30 {
		weight += 0.2
	}
	if event.Duration > 60 {
		weight += 0.3
	}

	// 根据行为类型调整权重
	switch event.BehaviorType {
	case "like":
		weight += 0.3
	case "comment":
		weight += 0.4
	case "share":
		weight += 0.5
	case "watch":
		// 观看权重已在时长中计算
	default:
		weight += 0.1
	}

	return weight
}

// updateRedisInterests 更新Redis中的用户兴趣标签
func (bp *BehaviorProcessor) updateRedisInterests(ctx context.Context, userID int64, videoTags []model.VideoTag, weightIncrement float64) error {
	for _, tag := range videoTags {
		// 计算最终权重增量（考虑标签权重）
		finalWeightDelta := weightIncrement * tag.Weight

		if err := redis.UpdateUserInterestWeight(ctx, userID, tag.TagName, finalWeightDelta); err != nil {
			klog.Errorf("Failed to update Redis interest weight: user_id=%d, tag=%s, error=%v",
				userID, tag.TagName, err)
		}
	}
	return nil
}

// updateMySQLInterests 更新MySQL中的用户兴趣标签
func (bp *BehaviorProcessor) updateMySQLInterests(ctx context.Context, event *rabbitmq.UserBehaviorEvent) error {
	// 获取视频标签
	var videoTags []model.VideoTag
	if err := db.DB.Where("video_id = ?", event.VideoID).Find(&videoTags).Error; err != nil {
		return fmt.Errorf("failed to get video tags: %v", err)
	}

	// 计算兴趣权重增量
	weightIncrement := bp.calculateWeightIncrement(event)

	// 更新用户兴趣标签
	for _, tag := range videoTags {
		var interest model.UserInterest
		result := db.DB.FirstOrCreate(&interest, model.UserInterest{
			UserID:  event.UserID,
			TagName: tag.TagName,
		})

		if result.Error != nil {
			klog.Errorf("Failed to find or create user interest: %v", result.Error)
			continue
		}

		// 使用指数衰减更新权重
		interest.Weight = interest.Weight*0.9 + weightIncrement*tag.Weight
		interest.UpdatedAt = time.Now()

		if err := db.DB.Save(&interest).Error; err != nil {
			klog.Errorf("Failed to save user interest: %v", err)
		}
	}

	return nil
}

// updateUserProfile 更新用户画像
func (bp *BehaviorProcessor) updateUserProfile(ctx context.Context, event *rabbitmq.UserBehaviorEvent) error {
	var profile model.UserProfile
	result := db.DB.FirstOrCreate(&profile, model.UserProfile{UserID: event.UserID})

	if result.Error != nil {
		return fmt.Errorf("failed to find or create user profile: %v", result.Error)
	}

	// 更新统计数据
	profile.LastActiveAt = time.Now()
	profile.TotalWatch += int64(event.Duration)

	// 根据行为类型更新相应计数
	switch event.BehaviorType {
	case "like":
		profile.TotalLikes++
	case "comment":
		profile.TotalComments++
	case "share":
		profile.TotalShares++
	}

	profile.UpdatedAt = time.Now()

	return db.DB.Save(&profile).Error
}

// BatchProcessInterests 批量处理用户兴趣标签更新
func (bp *BehaviorProcessor) BatchProcessInterests(ctx context.Context, updates []redis.UserInterestUpdate) error {
	// 批量更新Redis
	if err := redis.BatchUpdateUserInterests(ctx, updates); err != nil {
		return fmt.Errorf("failed to batch update Redis: %v", err)
	}

	// 批量更新MySQL
	return bp.batchUpdateMySQLInterests(ctx, updates)
}

// batchUpdateMySQLInterests 批量更新MySQL用户兴趣标签
func (bp *BehaviorProcessor) batchUpdateMySQLInterests(ctx context.Context, updates []redis.UserInterestUpdate) error {
	// 按用户ID分组
	userUpdates := make(map[int64][]redis.UserInterestUpdate)
	for _, update := range updates {
		userUpdates[update.UserID] = append(userUpdates[update.UserID], update)
	}

	// 批量处理每个用户
	for userID, userUpdateList := range userUpdates {
		if err := bp.processUserBatchUpdate(ctx, userID, userUpdateList); err != nil {
			klog.Errorf("Failed to process batch update for user %d: %v", userID, err)
		}
	}

	return nil
}

// processUserBatchUpdate 处理单个用户的批量更新
func (bp *BehaviorProcessor) processUserBatchUpdate(ctx context.Context, userID int64, updates []redis.UserInterestUpdate) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, update := range updates {
		var interest model.UserInterest
		result := tx.FirstOrCreate(&interest, model.UserInterest{
			UserID:  userID,
			TagName: update.TagName,
		})

		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find or create interest: %v", result.Error)
		}

		interest.Weight += update.WeightDelta
		interest.UpdatedAt = time.Now()

		if err := tx.Save(&interest).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save interest: %v", err)
		}
	}

	return tx.Commit().Error
}
