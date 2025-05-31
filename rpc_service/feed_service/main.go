package main

import (
	"log"
	feed "video_douyin/kitex_gen/feed/feedservice"
)

func main() {
	svr := feed.NewServer(new(FeedServiceImpl))
	if err := svr.Run(); err != nil {
		log.Fatal(err)
	}
}
