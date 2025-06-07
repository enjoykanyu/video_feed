package main

import (
	"log"
	recommend "video_douyin/kitex_gen/recommend/recommendservice"
)

func main() {
	svr := recommend.NewServer(new(RecommendServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
