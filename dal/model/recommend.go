package model

import (
	"time"
)

// UserInterest 用户兴趣标签表（优化版）
type UserInterest struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"index:idx_user_tag;not null" json:"user_id"`
	TagName     string    `gorm:"size:50;index:idx_user_tag;not null" json:"tag_name"`
	Weight      float64   `gorm:"type:decimal(10,4);default:0;not null" json:"weight"`
	DecayRate   float64   `gorm:"type:decimal(3,2);default:0.90" json:"decay_rate"`
	LastUpdate  time.Time `gorm:"index;not null" json:"last_update"`
	UpdateCount int32     `gorm:"default:0;not null" json:"update_count"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserInterest) TableName() string {
	return "user_interests"
}

// WatchHistory 用户观看历史表（优化版）
type WatchHistory struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"index:idx_user_video;not null" json:"user_id"`
	VideoID       int64     `gorm:"index:idx_user_video;not null" json:"video_id"`
	WatchDuration int32     `gorm:"not null;comment:观看时长(秒)" json:"watch_duration"`
	WatchProgress float64   `gorm:"type:decimal(5,2);default:0;comment:观看进度百分比" json:"watch_progress"`
	IsCompleted   bool      `gorm:"default:false;comment:是否完整观看" json:"is_completed"`
	DeviceType    string    `gorm:"size:20;comment:设备类型" json:"device_type"`
	IPAddress     string    `gorm:"size:45;comment:IP地址" json:"ip_address"`
	UserAgent     string    `gorm:"size:500;comment:用户代理" json:"user_agent"`
	CreatedAt     time.Time `gorm:"index;not null" json:"created_at"`

	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

func (WatchHistory) TableName() string {
	return "watch_histories"
}

// VideoTag 视频标签表（优化版）
type VideoTag struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VideoID    int64     `gorm:"index:idx_video_tag;not null" json:"video_id"`
	TagName    string    `gorm:"size:50;index:idx_video_tag;not null" json:"tag_name"`
	Weight     float64   `gorm:"type:decimal(10,4);default:1.0;not null" json:"weight"`
	TagType    string    `gorm:"size:20;default:'content';comment:标签类型" json:"tag_type"`
	Confidence float64   `gorm:"type:decimal(5,2);default:1.0;comment:置信度" json:"confidence"`
	Source     string    `gorm:"size:20;default:'manual';comment:标签来源" json:"source"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null" json:"updated_at"`

	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

func (VideoTag) TableName() string {
	return "video_tags"
}

// UserProfile 用户画像表（优化版）
type UserProfile struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"uniqueIndex;not null" json:"user_id"`
	LastActiveAt     time.Time `gorm:"index;not null" json:"last_active_at"`
	TotalWatch       int64     `gorm:"default:0;not null;comment:总观看时长(秒)" json:"total_watch"`
	TotalLikes       int64     `gorm:"default:0;not null;comment:总点赞数" json:"total_likes"`
	TotalComments    int64     `gorm:"default:0;not null;comment:总评论数" json:"total_comments"`
	TotalShares      int64     `gorm:"default:0;not null;comment:总分享数" json:"total_shares"`
	TotalFavorites   int64     `gorm:"default:0;not null;comment:总收藏数" json:"total_favorites"`
	AvgWatchDuration float64   `gorm:"type:decimal(8,2);default:0;comment:平均观看时长" json:"avg_watch_duration"`
	WatchFrequency   float64   `gorm:"type:decimal(5,2);default:0;comment:观看频率(次/天)" json:"watch_frequency"`
	PreferenceScore  float64   `gorm:"type:decimal(5,2);default:0;comment:偏好得分" json:"preference_score"`
	EngagementLevel  string    `gorm:"size:20;default:'low';comment:参与度等级" json:"engagement_level"`
	CreatedAt        time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"not null" json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

// VideoScore 视频评分表（优化版）
type VideoScore struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VideoID        int64     `gorm:"uniqueIndex;not null" json:"video_id"`
	HotScore       float64   `gorm:"type:decimal(10,4);default:0;not null" json:"hot_score"`
	QualityScore   float64   `gorm:"type:decimal(10,4);default:0;not null" json:"quality_score"`
	RelevanceScore float64   `gorm:"type:decimal(10,4);default:0;not null" json:"relevance_score"`
	FreshnessScore float64   `gorm:"type:decimal(10,4);default:0;not null" json:"freshness_score"`
	FinalScore     float64   `gorm:"type:decimal(10,4);default:0;not null" json:"final_score"`
	ScoreVersion   int32     `gorm:"default:1;not null;comment:评分版本" json:"score_version"`
	LastCalculated time.Time `gorm:"index;not null" json:"last_calculated"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`

	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

func (VideoScore) TableName() string {
	return "video_scores"
}

// UserBehaviorLog 用户行为日志表（新增）
type UserBehaviorLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"index:idx_user_time;not null" json:"user_id"`
	VideoID      int64     `gorm:"index:idx_video_time;not null" json:"video_id"`
	BehaviorType string    `gorm:"size:20;index;not null" json:"behavior_type"`
	Duration     int32     `gorm:"default:0;comment:行为持续时间" json:"duration"`
	Intensity    float64   `gorm:"type:decimal(5,2);default:1.0;comment:行为强度" json:"intensity"`
	Context      string    `gorm:"size:100;comment:行为上下文" json:"context"`
	DeviceInfo   string    `gorm:"size:200;comment:设备信息" json:"device_info"`
	Location     string    `gorm:"size:100;comment:地理位置" json:"location"`
	SessionID    string    `gorm:"size:64;comment:会话ID" json:"session_id"`
	CreatedAt    time.Time `gorm:"index;not null" json:"created_at"`

	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

func (UserBehaviorLog) TableName() string {
	return "user_behavior_logs"
}

// TagSimilarity 标签相似度表（新增）
type TagSimilarity struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TagName1     string    `gorm:"size:50;index:idx_tag_pair;not null" json:"tag_name1"`
	TagName2     string    `gorm:"size:50;index:idx_tag_pair;not null" json:"tag_name2"`
	Similarity   float64   `gorm:"type:decimal(5,4);default:0;not null" json:"similarity"`
	Correlation  float64   `gorm:"type:decimal(5,4);default:0;not null" json:"correlation"`
	Cooccurrence int64     `gorm:"default:0;not null;comment:共现次数" json:"cooccurrence"`
	LastUpdated  time.Time `gorm:"index;not null" json:"last_updated"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
}

func (TagSimilarity) TableName() string {
	return "tag_similarities"
}

// 常量定义
const (
	// 行为类型
	BehaviorTypeWatch    = "watch"
	BehaviorTypeLike     = "like"
	BehaviorTypeComment  = "comment"
	BehaviorTypeShare    = "share"
	BehaviorTypeFavorite = "favorite"

	// 标签类型
	TagTypeContent  = "content"
	TagTypeCategory = "category"
	TagTypeEmotion  = "emotion"

	// 标签来源
	TagSourceManual = "manual"
	TagSourceAI     = "ai"
	TagSourceAuto   = "auto"

	// 参与度等级
	EngagementLevelLow    = "low"
	EngagementLevelMedium = "medium"
	EngagementLevelHigh   = "high"
)
