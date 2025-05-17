package main

import (
	"log"
	user "video_douyin/kitex_gen/user/userservice"
)

func main() {
	//用 kitex 生成的代码创建一个 server 服务端
	svr := user.NewServer(new(UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
