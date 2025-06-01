namespace go upload

struct Tag {
    1: i64 id
    2: string name
    3: i64 created_at
}

struct UploadRequest {
    1: i64 user_id
    2: binary video_data
    3: string title
    4: string description
    5: list<string> tags
}

struct UploadResponse {
    1: bool success
    2: string message
    3: i64 video_id
}

struct UpdateRequest {
    1: i64 video_id
    2: string title
    3: string description
    4: list<string> tags
}

struct UpdateResponse {
    1: bool success
    2: string message
}

struct ProgressRequest {
    1: i64 upload_id
}

struct ProgressResponse {
    1: i32 progress
    2: string status
    3: string message
}

struct TagsRequest {
    1: string keyword
    2: i32 limit
}

struct TagsResponse {
    1: list<Tag> tags
}

service UploadService {
    UploadResponse UploadVideo(1: UploadRequest req)
    UpdateResponse UpdateVideoInfo(1: UpdateRequest req)
    ProgressResponse GetUploadProgress(1: ProgressRequest req)
    TagsResponse GetRecommendedTags(1: TagsRequest req)
} 