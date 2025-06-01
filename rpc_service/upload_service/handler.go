package main

import (
	"context"
	"time"
	"video_douyin/dal"
	"video_douyin/kitex_gen/upload"
	"video_douyin/pkg/minio"
)

type UploadServiceImpl struct{}

// UploadVideo 上传视频
func (s *UploadServiceImpl) UploadVideo(ctx context.Context, req *upload.UploadRequest) (*upload.UploadResponse, error) {
	// 生成唯一的视频文件名
	videoFileName := generateVideoFileName(req.UserId)

	// 保存视频到MinIO
	if err := minio.UploadVideo(videoFileName, string(req.VideoData)); err != nil {
		return nil, err
	}

	// 创建视频记录
	video := dal.Video{
		UserID:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
		VideoURL:    videoFileName,
		Status:      0, // 初始状态
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := dal.DB.Create(&video).Error; err != nil {
		return nil, err
	}

	// 处理标签
	for _, tagName := range req.Tags {
		var tag dal.Tag
		// 查找或创建标签
		dal.DB.FirstOrCreate(&tag, dal.Tag{Name: tagName})

		// 创建视频标签关联
		videoTag := dal.VideoTag{
			VideoID: video.ID,
			TagID:   tag.ID,
		}
		dal.DB.Create(&videoTag)
	}

	return &upload.UploadResponse{
		Success: true,
		Message: "Video uploaded successfully",
		VideoId: video.ID,
	}, nil
}

// UpdateVideoInfo 更新视频信息
func (s *UploadServiceImpl) UpdateVideoInfo(ctx context.Context, req *upload.UpdateRequest) (*upload.UpdateResponse, error) {
	// 更新视频基本信息
	if err := dal.DB.Model(&dal.Video{}).Where("id = ?", req.VideoId).
		Updates(map[string]interface{}{
			"title":       req.Title,
			"description": req.Description,
			"updated_at":  time.Now(),
		}).Error; err != nil {
		return nil, err
	}

	// 更新标签
	// 先删除旧的标签关联
	dal.DB.Where("video_id = ?", req.VideoId).Delete(&dal.VideoTag{})

	// 添加新的标签关联
	for _, tagName := range req.Tags {
		var tag dal.Tag
		dal.DB.FirstOrCreate(&tag, dal.Tag{Name: tagName})

		videoTag := dal.VideoTag{
			VideoID: req.VideoId,
			TagID:   tag.ID,
		}
		dal.DB.Create(&videoTag)
	}

	return &upload.UpdateResponse{
		Success: true,
		Message: "Video info updated successfully",
	}, nil
}

// GetUploadProgress 获取上传进度
func (s *UploadServiceImpl) GetUploadProgress(ctx context.Context, req *upload.ProgressRequest) (*upload.ProgressResponse, error) {
	// 这里可以实现一个简单的进度跟踪系统
	// 实际项目中可能需要使用Redis或其他存储来跟踪进度
	return &upload.ProgressResponse{
		Progress: 100, // 示例进度
		Status:   "completed",
		Message:  "Upload completed",
	}, nil
}

// GetRecommendedTags 获取推荐标签
func (s *UploadServiceImpl) GetRecommendedTags(ctx context.Context, req *upload.TagsRequest) (*upload.TagsResponse, error) {
	var tags []dal.Tag

	query := dal.DB.Model(&dal.Tag{})
	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	if err := query.Limit(int(req.Limit)).Find(&tags).Error; err != nil {
		return nil, err
	}

	respTags := make([]*upload.Tag, len(tags))
	for i, t := range tags {
		respTags[i] = &upload.Tag{
			Id:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt.Unix(),
		}
	}

	return &upload.TagsResponse{
		Tags: respTags,
	}, nil
}

// 生成视频文件名
func generateVideoFileName(userID int64) string {
	timestamp := time.Now().UnixNano()
	return string(userID) + "_" + string(timestamp) + ".mp4"
}
