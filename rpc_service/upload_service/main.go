package main

import (
	"log"
	upload "video_douyin/kitex_gen/upload/uploadservice"
)

func main() {
	svr := upload.NewServer(new(UploadServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
