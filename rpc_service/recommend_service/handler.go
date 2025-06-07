package main

import (
	"context"
	//"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	//"sort"
	"video_douyin/dal/model"
	//"encoding/json"
	"time"
	//"math"
	"github.com/go-redis/redis/v8"
	recommend "video_douyin/kitex_gen/recommend"
)

type RecommendServiceImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

// GetFeed 获取Feed流
func (s *RecommendServiceImpl) GetFeed(ctx context.Context, req *recommend.GetFeedRequest) (*recommend.GetFeedResponse, error) {
	// 1. 获取用户画像和兴趣标签
	userInterests, err := s.getUserInterests(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// 2. 获取候选视频列表
	candidateVideos, err := s.getCandidateVideos(ctx, req.UserId, req.LastTime, req.Count*3)
	if err != nil {
		return nil, err
	}

	// 3. 计算视频得分并排序
	scoredVideos := s.calculateVideoScores(candidateVideos, userInterests)

	// 4. 去重并返回结果
	finalVideos := s.deduplicateVideos(scoredVideos, req.Count)

	return &recommend.GetFeedResponse{
		Videos:   finalVideos,
		NextTime: time.Now().Unix(),
		HasMore:  len(finalVideos) == int(req.Count),
	}, nil
}

// getUserInterests 获取用户兴趣标签
func (s *RecommendServiceImpl) getUserInterests(ctx context.Context, userID int64) ([]*recommend.UserInterest, error) {
	// 先从Redis缓存获取
	//cacheKey := "user:interests:" + string(userID)
	//if interests, err := redis.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
	//	var result []*recommend.UserInterest
	//	if err := json.Unmarshal([]byte(interests), &result); err == nil {
	//		return result, nil
	//	}
	//}
	//
	//// 从数据库获取
	//var interests []model.UserInterest
	//if err := dal.DB.Where("user_id = ?", userID).Find(&interests).Error; err != nil {
	//	return nil, err
	//}
	//
	//// 转换为响应格式
	//result := make([]*recommend.UserInterest, len(interests))
	//for i, interest := range interests {
	//	result[i] = &recommend.UserInterest{
	//		UserId:    interest.UserID,
	//		TagName:   interest.TagName,
	//		Weight:    interest.Weight,
	//		UpdatedAt: interest.UpdatedAt.Unix(),
	//	}
	//}
	//
	//// 更新缓存
	//if data, err := json.Marshal(result); err == nil {
	//	redis.RedisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	//}

	return nil, nil
}

// getCandidateVideos 获取候选视频列表
func (s *RecommendServiceImpl) getCandidateVideos(ctx context.Context, userID int64, lastTime int64, count int32) ([]*model.Video, error) {
	//var videos []dal.Video
	//
	//// 基于用户兴趣标签获取相关视频
	//if err := dal.DB.Joins("JOIN video_tags ON videos.id = video_tags.video_id").
	//	Joins("JOIN user_interests ON video_tags.tag_name = user_interests.tag_name").
	//	Where("user_interests.user_id = ? AND videos.created_at < ?", userID, time.Unix(lastTime, 0)).
	//	Order("video_tags.weight * user_interests.weight DESC").
	//	Limit(int(count)).
	//	Find(&videos).Error; err != nil {
	//	return nil, err
	//}
	//
	//return &videos, nil
	return nil, nil
}

// calculateVideoScores 计算视频得分
func (s *RecommendServiceImpl) calculateVideoScores(videos []*model.Video, userInterests []*recommend.UserInterest) []*recommend.Video {
	//scoredVideos := make([]*recommend.Video, len(videos))
	//
	//for i, video := range videos {
	//	// 获取视频标签
	//	var videoTags []dal.VideoTag
	//	dal.DB.Where("video_id = ?", video.ID).Find(&videoTags)
	//
	//	// 计算相关性得分
	//	relevanceScore := s.calculateRelevanceScore(videoTags, userInterests)
	//
	//	// 获取视频热度得分
	//	var videoScore dal.VideoScore
	//	dal.DB.Where("video_id = ?", video.ID).First(&videoScore)
	//
	//	// 综合得分 = 相关性得分 * 0.6 + 热度得分 * 0.4
	//	finalScore := relevanceScore*0.6 + videoScore.HotScore*0.4
	//
	//	scoredVideos[i] = &recommend.Video{
	//		Id:          video.ID,
	//		UserId:      video.UserID,
	//		Title:       video.Title,
	//		Description: video.Description,
	//		CoverUrl:    video.CoverURL,
	//		VideoUrl:    video.VideoURL,
	//		Status:      video.Status,
	//		CreatedAt:   video.CreatedAt.Unix(),
	//		UpdatedAt:   video.UpdatedAt.Unix(),
	//		Score:       finalScore,
	//	}
	//}
	//
	//// 按得分排序
	//sort.Slice(scoredVideos, func(i, j int) bool {
	//	return scoredVideos[i].Score > scoredVideos[j].Score
	//})
	//
	//return scoredVideos
	return nil
}

// calculateRelevanceScore 计算相关性得分
func (s *RecommendServiceImpl) calculateRelevanceScore(videoTags []model.VideoTag, userInterests []*recommend.UserInterest) float64 {
	var score float64

	//// 计算标签匹配度
	//for _, videoTag := range videoTags {
	//	for _, userInterest := range userInterests {
	//		if videoTag.TagName == userInterest.TagName {
	//			score += videoTag.Weight * userInterest.Weight
	//		}
	//	}
	//}

	return score
}

// deduplicateVideos 视频去重
func (s *RecommendServiceImpl) deduplicateVideos(videos []*recommend.Video, count int32) []*recommend.Video {
	// 使用Redis记录已推荐视频
	//ctx := context.Background()
	//recommendedKey := "user:recommended:" + string(videos[0].UserId)
	//
	//// 获取已推荐视频列表
	//recommended, _ := redis.RedisClient.SMembers(ctx, recommendedKey).Result()
	//recommendedMap := make(map[string]bool)
	//for _, v := range recommended {
	//	recommendedMap[v] = true
	//}
	//
	//// 过滤并返回未推荐视频
	//result := make([]*recommend.Video, 0, count)
	//for _, video := range videos {
	//	if len(result) >= int(count) {
	//		break
	//	}
	//
	//	videoID := string(video.Id)
	//	if !recommendedMap[videoID] {
	//		result = append(result, video)
	//		// 记录已推荐视频
	//		redis.RedisClient.SAdd(ctx, recommendedKey, videoID)
	//	}
	//}
	//
	//// 设置过期时间（24小时）
	//redis.RedisClient.Expire(ctx, recommendedKey, 24*time.Hour)

	return nil
}

// UpdateUserProfile 更新用户画像
func (s *RecommendServiceImpl) UpdateUserProfile(ctx context.Context, req *recommend.UpdateUserProfileRequest) (*recommend.UpdateUserProfileResponse, error) {
	// 1. 更新观看历史
	//watchHistory := dal.WatchHistory{
	//	UserID:        req.UserId,
	//	VideoID:       req.VideoId,
	//	WatchDuration: req.WatchDuration,
	//	CreatedAt:     time.Now(),
	//}
	//if err := dal.DB.Create(&watchHistory).Error; err != nil {
	//	return nil, err
	//}
	//
	//// 2. 更新用户画像
	//var profile dal.UserProfile
	//dal.DB.FirstOrCreate(&profile, dal.UserProfile{UserID: req.UserId})
	//
	//profile.LastActiveAt = time.Now()
	//profile.TotalWatch += int64(req.WatchDuration)
	//if req.IsLike {
	//	profile.TotalLikes++
	//}
	//if req.IsComment {
	//	profile.TotalComments++
	//}
	//if req.IsShare {
	//	profile.TotalShares++
	//}
	//
	//if err := dal.DB.Save(&profile).Error; err != nil {
	//	return nil, err
	//}
	//
	//// 3. 更新用户兴趣标签
	//if err := s.updateUserInterests(ctx, req); err != nil {
	//	return nil, err
	//}
	//
	//return &recommend.UpdateUserProfileResponse{
	//	Success: true,
	//	Message: "User profile updated successfully",
	//}, nil
	return nil, nil
}

// updateUserInterests 更新用户兴趣标签
func (s *RecommendServiceImpl) updateUserInterests(ctx context.Context, req *recommend.UpdateUserProfileRequest) error {
	// 获取视频标签
	//var videoTags []dal.VideoTag
	//if err := dal.DB.Where("video_id = ?", req.VideoId).Find(&videoTags).Error; err != nil {
	//	return err
	//}
	//
	//// 计算兴趣权重增量
	//weightIncrement := s.calculateInterestWeight(req)
	//
	//// 更新用户兴趣标签
	//for _, tag := range videoTags {
	//	var interest dal.UserInterest
	//	result := dal.DB.FirstOrCreate(&interest, dal.UserInterest{
	//		UserID:  req.UserId,
	//		TagName: tag.TagName,
	//	})
	//
	//	if result.Error != nil {
	//		continue
	//	}
	//
	//	// 使用指数衰减更新权重
	//	interest.Weight = interest.Weight*0.9 + weightIncrement*tag.Weight
	//	interest.UpdatedAt = time.Now()
	//
	//	if err := dal.DB.Save(&interest).Error; err != nil {
	//		klog.Errorf("Failed to update user interest: %v", err)
	//	}
	//}

	return nil
}

// calculateInterestWeight 计算兴趣权重增量
func (s *RecommendServiceImpl) calculateInterestWeight(req *recommend.UpdateUserProfileRequest) float64 {
	var weight float64 = 0.1 // 基础权重

	// 根据观看时长调整权重
	if req.WatchDuration > 30 {
		weight += 0.2
	}

	// 根据互动行为调整权重
	if req.IsLike {
		weight += 0.3
	}
	if req.IsComment {
		weight += 0.4
	}
	if req.IsShare {
		weight += 0.5
	}

	return weight
}
