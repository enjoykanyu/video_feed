package main

import (
	"context"
	"encoding/json"
	"time"
	"video_douyin/dal"
	"video_douyin/kitex_gen/interaction"

	// InteractionService "video_douyin/kitex_gen/interaction"
	"video_douyin/pkg/redis"
)

type InteractionServiceImpl struct{}

// Follow 关注用户
func (s *InteractionServiceImpl) Follow(ctx context.Context, req *interaction.FollowRequest) (*interaction.FollowResponse, error) {
	// 创建关注记录 //局部变量隐蔽 interaction和返回包定义相同 所以会报错interaction.FollowResponse is not a typecompilerNotAType
	interaction_db := dal.Interaction{
		UserID:    req.UserId,
		VideoID:   req.TargetId,
		Type:      3, // 关注类型
		CreatedAt: time.Now(),
	}

	if err := dal.DB.Create(&interaction_db).Error; err != nil {
		return nil, err
	}

	// 更新Redis缓存
	following, _ := redis.GetFollowing(ctx, req.UserId)
	var followingList []int64
	if following != "" {
		json.Unmarshal([]byte(following), &followingList)
	}
	followingList = append(followingList, req.TargetId)
	followingData, _ := json.Marshal(followingList)
	redis.CacheFollowing(ctx, req.UserId, string(followingData))

	return &interaction.FollowResponse{
		Success: true,
		Message: "Followed successfully",
	}, nil
}

// Unfollow 取消关注
func (s *InteractionServiceImpl) Unfollow(ctx context.Context, req *interaction.FollowRequest) (*interaction.FollowResponse, error) {
	// 删除关注记录
	if err := dal.DB.Where("user_id = ? AND video_id = ? AND type = 3", req.UserId, req.TargetId).Delete(&dal.Interaction{}).Error; err != nil {
		return nil, err
	}

	// 更新Redis缓存
	following, _ := redis.GetFollowing(ctx, req.UserId)
	var followingList []int64
	if following != "" {
		json.Unmarshal([]byte(following), &followingList)
		// 从列表中移除目标用户
		for i, id := range followingList {
			if id == req.TargetId {
				followingList = append(followingList[:i], followingList[i+1:]...)
				break
			}
		}
		followingData, _ := json.Marshal(followingList)
		redis.CacheFollowing(ctx, req.UserId, string(followingData))
	}

	return &interaction.FollowResponse{
		Success: true,
		Message: "Unfollowed successfully",
	}, nil
}

// Comment 发表评论
func (s *InteractionServiceImpl) Comment(ctx context.Context, req *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	comment := dal.Comment{
		UserID:    req.UserId,
		VideoID:   req.VideoId,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := dal.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &interaction.CommentResponse{
		Success: true,
		Message: "Comment posted successfully",
		Comment: &interaction.Comment{
			Id:        comment.ID,
			UserId:    comment.UserID,
			VideoId:   comment.VideoID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Unix(),
		},
	}, nil
}

// GetComments 获取评论列表
func (s *InteractionServiceImpl) GetComments(ctx context.Context, req *interaction.GetCommentsRequest) (*interaction.GetCommentsResponse, error) {
	var comments []dal.Comment
	offset := (req.Page - 1) * req.PageSize

	if err := dal.DB.Where("video_id = ?", req.VideoId).
		Order("created_at desc").
		Offset(int(offset)).
		Limit(int(req.PageSize)).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	var total int64
	dal.DB.Model(&dal.Comment{}).Where("video_id = ?", req.VideoId).Count(&total)

	respComments := make([]*interaction.Comment, len(comments))
	for i, c := range comments {
		respComments[i] = &interaction.Comment{
			Id:        c.ID,
			UserId:    c.UserID,
			VideoId:   c.VideoID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Unix(),
		}
	}

	return &interaction.GetCommentsResponse{
		Comments: respComments,
		Total:    int32(total),
	}, nil
}

// Like 点赞
func (s *InteractionServiceImpl) Like(ctx context.Context, req *interaction.LikeRequest) (*interaction.LikeResponse, error) {
	interaction_db := dal.Interaction{
		UserID:    req.UserId,
		VideoID:   req.VideoId,
		Type:      1, // 点赞类型
		CreatedAt: time.Now(),
	}

	if err := dal.DB.Create(&interaction_db).Error; err != nil {
		return nil, err
	}

	// 更新Redis缓存
	interactions, _ := redis.GetUserInteractions(ctx, req.UserId)
	var interactionList []int64
	if interactions != "" {
		json.Unmarshal([]byte(interactions), &interactionList)
	}
	interactionList = append(interactionList, req.VideoId)
	interactionData, _ := json.Marshal(interactionList)
	redis.CacheUserInteractions(ctx, req.UserId, string(interactionData))

	return &interaction.LikeResponse{
		Success: true,
		Message: "Liked successfully",
	}, nil
}

// Unlike 取消点赞
func (s *InteractionServiceImpl) Unlike(ctx context.Context, req *interaction.LikeRequest) (*interaction.LikeResponse, error) {
	if err := dal.DB.Where("user_id = ? AND video_id = ? AND type = 1", req.UserId, req.VideoId).Delete(&dal.Interaction{}).Error; err != nil {
		return nil, err
	}

	// 更新Redis缓存
	interactions, _ := redis.GetUserInteractions(ctx, req.UserId)
	var interactionList []int64
	if interactions != "" {
		json.Unmarshal([]byte(interactions), &interactionList)
		// 从列表中移除视频ID
		for i, id := range interactionList {
			if id == req.VideoId {
				interactionList = append(interactionList[:i], interactionList[i+1:]...)
				break
			}
		}
		interactionData, _ := json.Marshal(interactionList)
		redis.CacheUserInteractions(ctx, req.UserId, string(interactionData))
	}

	return &interaction.LikeResponse{
		Success: true,
		Message: "Unliked successfully",
	}, nil
}

// Favorite 收藏
func (s *InteractionServiceImpl) Favorite(ctx context.Context, req *interaction.FavoriteRequest) (*interaction.FavoriteResponse, error) {
	interaction_db := dal.Interaction{
		UserID:    req.UserId,
		VideoID:   req.VideoId,
		Type:      2, // 收藏类型
		CreatedAt: time.Now(),
	}

	if err := dal.DB.Create(&interaction_db).Error; err != nil {
		return nil, err
	}

	return &interaction.FavoriteResponse{
		Success: true,
		Message: "Favorited successfully",
	}, nil
}

// Unfavorite 取消收藏
func (s *InteractionServiceImpl) Unfavorite(ctx context.Context, req *interaction.FavoriteRequest) (*interaction.FavoriteResponse, error) {
	if err := dal.DB.Where("user_id = ? AND video_id = ? AND type = 2", req.UserId, req.VideoId).Delete(&dal.Interaction{}).Error; err != nil {
		return nil, err
	}

	return &interaction.FavoriteResponse{
		Success: true,
		Message: "Unfavorited successfully",
	}, nil
}

// SendDanmaku 发送弹幕
func (s *InteractionServiceImpl) SendDanmaku(ctx context.Context, req *interaction.DanmakuRequest) (*interaction.DanmakuResponse, error) {
	danmaku := dal.Danmaku{
		UserID:    req.UserId,
		VideoID:   req.VideoId,
		Content:   req.Content,
		Position:  req.Position,
		Color:     req.Color,
		CreatedAt: time.Now(),
	}

	if err := dal.DB.Create(&danmaku).Error; err != nil {
		return nil, err
	}

	return &interaction.DanmakuResponse{
		Success: true,
		Message: "Danmaku sent successfully",
	}, nil
}

// GetDanmaku 获取弹幕列表
func (s *InteractionServiceImpl) GetDanmaku(ctx context.Context, req *interaction.GetDanmakuRequest) (*interaction.GetDanmakuResponse, error) {
	var danmakuList []dal.Danmaku
	if err := dal.DB.Where("video_id = ? AND position BETWEEN ? AND ?", req.VideoId, req.StartTime, req.EndTime).
		Order("position asc").
		Find(&danmakuList).Error; err != nil {
		return nil, err
	}

	respDanmakuList := make([]*interaction.Danmaku, len(danmakuList))
	for i, d := range danmakuList {
		respDanmakuList[i] = &interaction.Danmaku{
			Id:        d.ID,
			UserId:    d.UserID,
			VideoId:   d.VideoID,
			Content:   d.Content,
			Position:  d.Position,
			Color:     d.Color,
			CreatedAt: d.CreatedAt.Unix(),
		}
	}

	return &interaction.GetDanmakuResponse{
		DanmakuList: respDanmakuList,
	}, nil
}

// Share 分享视频
func (s *InteractionServiceImpl) Share(ctx context.Context, req *interaction.ShareRequest) (*interaction.ShareResponse, error) {
	// 生成分享链接
	shareURL := generateShareURL(req.VideoId, req.Platform)

	return &interaction.ShareResponse{
		Success:  true,
		ShareUrl: shareURL,
		Message:  "Share URL generated successfully",
	}, nil
}

// 生成分享链接
func generateShareURL(videoID int64, platform string) string {
	// 这里可以根据不同平台生成不同的分享链接
	baseURL := "https://douyin.com/video/"
	return baseURL + string(videoID)
}
