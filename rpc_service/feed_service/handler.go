package main

import (
	"context"
	"encoding/json"
	"time"
	"video_douyin/dal"
	"video_douyin/kitex_gen/feed"
	"video_douyin/pkg/minio"
	"video_douyin/pkg/redis"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FeedServiceImpl struct{}

// GetFeed 获取视频Feed流
func (s *FeedServiceImpl) GetFeed(ctx context.Context, req *feed.FeedRequest) (*feed.FeedResponse, error) {
	// 尝试从Redis缓存获取
	if req.LastTime == 0 {
		hotVideos, err := redis.GetHotVideos(ctx)
		if err == nil {
			var videos []*feed.Video
			if err := json.Unmarshal([]byte(hotVideos), &videos); err == nil {
				return &feed.FeedResponse{
					Videos:   videos,
					NextTime: time.Now().Unix(),
					HasMore:  true,
				}, nil
			}
		}
	}

	// 从数据库获取视频列表
	var videos []dal.Video
	err := dal.DB.Where("created_at < ?", time.Unix(req.LastTime, 0)).
		Order("created_at desc").
		Limit(int(req.Count)).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	respVideos := make([]*feed.Video, len(videos))
	for i, v := range videos {
		respVideos[i] = &feed.Video{
			Id:          v.ID,
			UserId:      v.UserID,
			Title:       v.Title,
			Description: v.Description,
			CoverUrl:    v.CoverURL,
			VideoUrl:    v.VideoURL,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt.Unix(),
			UpdatedAt:   v.UpdatedAt.Unix(),
		}
	}

	return &feed.FeedResponse{
		Videos:   respVideos,
		NextTime: time.Now().Unix(),
		HasMore:  len(videos) == int(req.Count),
	}, nil
}

// GetVideoChunk 获取视频分片
func (s *FeedServiceImpl) GetVideoChunk(ctx context.Context, req *feed.ChunkRequest) (*feed.ChunkResponse, error) {
	// 获取视频信息
	var video dal.Video
	if err := dal.DB.First(&video, req.VideoId).Error; err != nil {
		return nil, err
	}

	// 从MinIO获取视频分片
	chunkSize := int64(1024 * 1024) // 1MB per chunk
	offset := int64(req.ChunkIndex) * chunkSize
	data, err := minio.GetVideoChunk(video.VideoURL, offset, chunkSize)
	if err != nil {
		return nil, err
	}

	return &feed.ChunkResponse{
		Data:         data,
		TotalChunks:  int32(video.Status), // 使用Status字段存储总分片数
		CurrentChunk: req.ChunkIndex,
	}, nil
}

// PreloadVideo 预加载视频
func (s *FeedServiceImpl) PreloadVideo(ctx context.Context, req *feed.PreloadRequest) (*feed.PreloadResponse, error) {
	// 获取视频信息
	var video dal.Video
	if err := dal.DB.First(&video, req.VideoId).Error; err != nil {
		return nil, err
	}

	// 根据优先级决定预加载策略
	switch req.Priority {
	case 1: // 高优先级，立即预加载
		go func() {
			// 预加载前几个分片
			for i := 0; i < 3; i++ {
				_, err := minio.GetVideoChunk(video.VideoURL, int64(i)*1024*1024, 1024*1024)
				if err != nil {
					klog.Errorf("Failed to preload video chunk: %v", err)
					return
				}
			}
		}()
	case 2: // 中优先级，延迟预加载
		time.AfterFunc(5*time.Second, func() {
			// 预加载前两个分片
			for i := 0; i < 2; i++ {
				_, err := minio.GetVideoChunk(video.VideoURL, int64(i)*1024*1024, 1024*1024)
				if err != nil {
					klog.Errorf("Failed to preload video chunk: %v", err)
					return
				}
			}
		})
	case 3: // 低优先级，仅记录
		klog.Infof("Video %d marked for low priority preload", req.VideoId)
	}

	return &feed.PreloadResponse{
		Success: true,
		Message: "Preload scheduled",
	}, nil
}
