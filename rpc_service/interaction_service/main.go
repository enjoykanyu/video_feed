package main

import (
	"log"
	interaction "video_douyin/kitex_gen/interaction/interactionservice"
)

func main() {
	svr := interaction.NewServer(new(InteractionServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
