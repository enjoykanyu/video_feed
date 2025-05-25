package main

import (
	"log"
	user "video_douyin/kitex_gen/user/userservice"
	"video_douyin/pkg/db"
)

func main() {
	// 初始化MySQL
	if err := db.InitMySQL("root:901project@tcp(127.0.0.1:3306)/kanyuServer?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
		log.Fatalf("MySQL初始化失败: %v", err)
	}

	// 初始化Redis
	if err := db.InitRedis("127.0.0.1:6379", 0); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}

	// 创建UserServiceImpl实例并设置db和redis
	impl := &UserServiceImpl{
		db:    db.DB,
		redis: db.Redis,
	}

	// 创建server
	svr := user.NewServer(impl)

	// 运行服务
	if err := svr.Run(); err != nil {
		log.Fatalf("服务运行失败: %v", err)
	}
}
