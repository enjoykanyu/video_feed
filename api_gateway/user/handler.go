package user

import (
	"context"
	"log"
	"net/http"
	"time"
	"video_douyin/kitex_gen/user"             // Kitex生成的用户服务数据结构
	"video_douyin/kitex_gen/user/userservice" // Kitex生成的用户服务客户端

	"github.com/cloudwego/hertz/pkg/app" // Hertz HTTP上下文处理
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/client" // Kitex客户端核心
	"github.com/cloudwego/kitex/client/callopt"
	// Kitex调用选项
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

// type RegisterRequest struct {
// 	Phone      string `thrift:"phone,1,required" frugal:"1,required,string" json:"phone"`
// 	VerifyCode string `thrift:"verifyCode,2,required" frugal:"2,required,string" json:"verifyCode"`
// }

// func NewRegisterRequest() *RegisterRequest {
// 	return &RegisterRequest{}
// }

// Handler 处理用户注册HTTP请求
func Register(ctx context.Context, c *app.RequestContext) {
	req := user.NewRegisterRequest()
	// var RegisterRequest struct {
	// 	Phone      string `thrift:"phone,1,required" frugal:"1,required,string" json:"phone"`
	// 	VerifyCode string `thrift:"verifyCode,2,required" frugal:"2,required,string" json:"verifyCode"`
	// }

	// req := NewRegisterRequest()

	if err := c.BindAndValidate(req); err != nil {
		c.JSON(http.StatusOK, utils.H{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}
	log.Printf("请求参数: %v\n", req.Phone)      // 打印请求参数
	log.Printf("请求参数: %v\n", req.VerifyCode) // 打印请求参数
	log.Printf("参数类型: %T\n", req)            // 打印请求参数
	// 调用用户服务RPC接口，超时3秒
	resp, err := userClient.Register(
		ctx,
		req,
		callopt.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatal(err) // RPC调用失败记录日志
	}
	c.String(200, resp.String()) // 返回RPC响应内容
}

func Login(ctx context.Context, c *app.RequestContext) {
	req := user.NewLoginRequest()
	// var RegisterRequest struct {
	// 	Phone      string `thrift:"phone,1,required" frugal:"1,required,string" json:"phone"`
	// 	VerifyCode string `thrift:"verifyCode,2,required" frugal:"2,required,string" json:"verifyCode"`
	// }

	// req := NewRegisterRequest()

	if err := c.BindAndValidate(req); err != nil {
		c.JSON(http.StatusOK, utils.H{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}
	log.Printf("请求参数: %v\n", req.Phone)         // 打印请求参数
	log.Printf("请求参数验证码: %v\n", req.VerifyCode) // 打印请求参数
	log.Printf("参数类型: %T\n", req)               // 打印请求参数
	// 调用用户服务RPC接口，超时3秒
	resp, err := userClient.Login(
		ctx,
		req,
		callopt.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatal(err) // RPC调用失败记录日志
	}
	//c.String(200, resp.String()) // 返回RPC响应内容
	c.JSON(200, resp)
}
func Code(ctx context.Context, c *app.RequestContext) {
	req := user.NewSendVerifyCodeRequest()
	phone := c.Query("phone")
	req.Phone = phone
	log.Printf("请求参数: %v\n", phone) // 打印请求参数
	log.Printf("参数类型: %T\n", req)   // 打印请求参数
	// 调用用户服务RPC接口，超时3秒
	resp, err := userClient.SendVerifyCode(
		ctx,
		req,
		callopt.WithRPCTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatal(err) // RPC调用失败记录日志
	}
	c.String(200, resp.String()) // 返回RPC响应内容
}
