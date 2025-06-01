namespace go interaction

struct Comment {
    1: i64 id
    2: i64 user_id
    3: i64 video_id
    4: string content
    5: i64 created_at
}

struct Danmaku {
    1: i64 id
    2: i64 user_id
    3: i64 video_id
    4: string content
    5: i32 position
    6: string color
    7: i64 created_at
}

struct FollowRequest {
    1: i64 user_id
    2: i64 target_id
}

struct FollowResponse {
    1: bool success
    2: string message
}

struct CommentRequest {
    1: i64 user_id
    2: i64 video_id
    3: string content
}

struct CommentResponse {
    1: bool success
    2: string message
    3: Comment comment
}

struct GetCommentsRequest {
    1: i64 video_id
    2: i32 page
    3: i32 page_size
}

struct GetCommentsResponse {
    1: list<Comment> comments
    2: i32 total
}

struct LikeRequest {
    1: i64 user_id
    2: i64 video_id
}

struct LikeResponse {
    1: bool success
    2: string message
}

struct FavoriteRequest {
    1: i64 user_id
    2: i64 video_id
}

struct FavoriteResponse {
    1: bool success
    2: string message
}

struct DanmakuRequest {
    1: i64 user_id
    2: i64 video_id
    3: string content
    4: i32 position
    5: string color
}

struct DanmakuResponse {
    1: bool success
    2: string message
}

struct GetDanmakuRequest {
    1: i64 video_id
    2: i32 start_time
    3: i32 end_time
}

struct GetDanmakuResponse {
    1: list<Danmaku> danmaku_list
}

struct ShareRequest {
    1: i64 user_id
    2: i64 video_id
    3: string platform
}

struct ShareResponse {
    1: bool success
    2: string share_url
    3: string message
}

service InteractionService {
    FollowResponse Follow(1: FollowRequest req)
    FollowResponse Unfollow(1: FollowRequest req)
    
    CommentResponse Comment(1: CommentRequest req)
    GetCommentsResponse GetComments(1: GetCommentsRequest req)
    
    LikeResponse Like(1: LikeRequest req)
    LikeResponse Unlike(1: LikeRequest req)
    
    FavoriteResponse Favorite(1: FavoriteRequest req)
    FavoriteResponse Unfavorite(1: FavoriteRequest req)
    
    DanmakuResponse SendDanmaku(1: DanmakuRequest req)
    GetDanmakuResponse GetDanmaku(1: GetDanmakuRequest req)
    
    ShareResponse Share(1: ShareRequest req)
} 