namespace go recommend

// 用户兴趣标签
struct UserInterest {
    1: i64 user_id
    2: string tag_name
    3: double weight        // 兴趣权重
    4: i64 updated_at
}

// 用户观看历史
struct WatchHistory {
    1: i64 user_id
    2: i64 video_id
    3: i32 watch_duration  // 观看时长(秒)
    4: i64 created_at
}
struct Video {
    1: i64 id
    2: i64 user_id
    3: string title
    4: string description
    5: string cover_url
    6: string video_url
    7: i32 status
    8: i64 created_at
    9: i64 updated_at
    10: double  score
}

// 视频标签
struct VideoTag {
    1: i64 video_id
    2: string tag_name
    3: double weight       // 标签权重
}

// 获取Feed请求
struct GetFeedRequest {
    1: i64 user_id
    2: i64 last_time      // 上次请求时间
    3: i32 count          // 请求数量
    4: bool is_refresh    // 是否刷新
}

// 获取Feed响应
struct GetFeedResponse {
    1: list<Video> videos
    2: i64 next_time
    3: bool has_more
}

// 更新用户画像请求
struct UpdateUserProfileRequest {
    1: i64 user_id
    2: i64 video_id
    3: i32 watch_duration
    4: bool is_like
    5: bool is_comment
    6: bool is_share
}

// 更新用户画像响应
struct UpdateUserProfileResponse {
    1: bool success
    2: string message
}

// 获取用户兴趣标签请求
struct GetUserInterestsRequest {
    1: i64 user_id
}

// 获取用户兴趣标签响应
struct GetUserInterestsResponse {
    1: list<UserInterest> interests
}

service RecommendService {
    // 获取Feed流
    GetFeedResponse GetFeed(1: GetFeedRequest req)

    // 更新用户画像
    UpdateUserProfileResponse UpdateUserProfile(1: UpdateUserProfileRequest req)

    // 获取用户兴趣标签
    GetUserInterestsResponse GetUserInterests(1: GetUserInterestsRequest req)
}