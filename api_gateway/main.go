package main

import (
	"log"
	"video_douyin/middleware"

	"api_gateway/user"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 初始化 Hertz 服务器实例，监听本地 8889 端口
	// WithHostPorts 配置项指定服务监听地址:ml-citation{ref="6,7" data="citationList"}
	hz := server.New(server.WithHostPorts("localhost:8889"))
	// 注册全局中间件
	hz.Use(middleware.AuthMiddleware())
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
