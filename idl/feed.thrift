namespace go feed

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
}

struct FeedRequest {
    1: i64 user_id
    2: i64 last_time
    3: i32 count
}

struct FeedResponse {
    1: list<Video> videos
    2: i64 next_time
    3: bool has_more
}

struct ChunkRequest {
    1: i64 video_id
    2: i32 chunk_index
}

struct ChunkResponse {
    1: binary data
    2: i32 total_chunks
    3: i32 current_chunk
}

struct PreloadRequest {
    1: i64 video_id
    2: i32 priority
}

struct PreloadResponse {
    1: bool success
    2: string message
}

service FeedService {
    FeedResponse GetFeed(1: FeedRequest req)
    ChunkResponse GetVideoChunk(1: ChunkRequest req)
    PreloadResponse PreloadVideo(1: PreloadRequest req)
} 