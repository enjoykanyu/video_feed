package main

import (
	"log"

	"api_gateway/user"
	// "video_douyin/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// var (
// 	cli userservice.Client
// )

func main() {
	// // 初始化MySQL
	// if err := db.InitMySQL("root:901project@tcp(127.0.0.1:3306)/kanyuServer?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
	// 	log.Fatalf("MySQL初始化失败: %v", err)
	// }
	// // c, err := userservice.NewClient("douyin.video.user", client.WithHostPorts("0.0.0.0:8888"))
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // cli = c

	// // 初始化Redis
	// db.InitRedis("127.0.0.1:6379", "", 0)
	// 初始化 Hertz 服务器实例，监听本地 8889 端口
	// WithHostPorts 配置项指定服务监听地址:ml-citation{ref="6,7" data="citationList"}
	hz := server.New(server.WithHostPorts("localhost:8889"))

	// 创建路由分组（符合 RESTful 风格）
	// 第一级分组 /douyin 作为 API 根路径:ml-citation{ref="2" data="citationList"}
	douyin := hz.Group("/douyin")
	// 第二级分组 /user 用于用户相关操作:ml-citation{ref="3" data="citationList"}
	userGroup := douyin.Group("/user")
	// 注册用户注册接口，POST 方法对应创建操作
	// 处理函数指向 user 包的 Handler 方法:ml-citation{ref="3" data="citationList"}
	userGroup.POST("/register", user.Register)
	userGroup.POST("/login/code", user.Login)
	userGroup.POST("/code", user.Code)
	// 启动 HTTP 服务，若失败则记录错误日志
	// Run() 会阻塞直到服务终止:ml-citation{ref="6" data="citationList"}
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// func Handler(ctx context.Context, c *app.RequestContext) {
// req := user.NewGetUserInfoRequest()
// resp, err := cli.GetUserInfo(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
// if err != nil {
// 	log.Fatal(err)
// }

// c.String(200, resp.String())
// }
