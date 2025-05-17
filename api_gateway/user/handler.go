package user

import (
	"context"
	"log"
	"time"
	"video_douyin/kitex_gen/user"             // Kitex生成的用户服务数据结构
	"video_douyin/kitex_gen/user/userservice" // Kitex生成的用户服务客户端

	"github.com/cloudwego/hertz/pkg/app"        // Hertz HTTP上下文处理
	"github.com/cloudwego/kitex/client"         // Kitex客户端核心
	"github.com/cloudwego/kitex/client/callopt" // Kitex调用选项
)

// 全局用户服务客户端实例
var userClient userservice.Client

// 初始化Kitex客户端连接
func init() {
	c, err := userservice.NewClient(
		"douyin.video.user",                  // 服务名
		client.WithHostPorts("0.0.0.0:8888"), // 服务地址
	)
	if err != nil {
		panic(err) // 连接失败终止程序
	}
	userClient = c
}

// Handler 处理用户注册HTTP请求
func Handler(ctx context.Context, c *app.RequestContext) {
	req := user.NewRegisterRequest()   // 创建注册请求结构体
	log.Printf("Phone: %s", req.Phone) // 记录请求参数

	// 调用用户服务RPC接口，超时3秒
	resp, err := userClient.Register(
		context.Background(),
		req,
		callopt.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatal(err) // RPC调用失败记录日志
	}
	c.String(200, resp.String()) // 返回RPC响应内容
}
