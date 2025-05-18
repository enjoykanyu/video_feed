package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	user "video_douyin/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// TODO: Your code here...

	resp = user.NewRegisterResponse()

	// 1. Validate the request parameters
	if req.GetPhone() == "" || req.GetVerifyCode() == "" {
		errorMsg := "Invalid request parameters" // 声明字符串变量
		resp.SetMessage(&errorMsg)               // 传递指针
		resp.SetSuccess(true)
		return
	}

	// 2. Check if the username already exists
	// ...

	// 3. Create a new user
	// ...

	// 4. Return the response

	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	// TODO: Your code here...
	resp = user.NewLoginResponse()
	// 1. 校验参数异常
	if req.GetPhone() == "" || req.GetVerifyCode() == "" {

		errorMsg := "Invalid request parameters" // 声明字符串变量
		resp.SetMessage(&errorMsg)               // 传递指针
		resp.SetSuccess(true)
	}
	// 2. 查询db

	// 3. 生成token
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// SendVerifyCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) SendVerifyCode(ctx context.Context, req *user.SendVerifyCodeRequest) (resp *user.SendVerifyCodeResponse, err error) {
	// TODO: Your code here...
	resp = user.NewSendVerifyCodeResponse()
	// 1. 校验参数异常
	if req.GetPhone() == "" {

		errorMsg := "Invalid request parameters" // 声明字符串变量
		resp.SetMessage(&errorMsg)               // 传递指针
		resp.SetSuccess(false)
	}
	//2,查询db看用户手机号是否存在

	//3, 生成验证码
	code, _ := generateCaptcha()

	log.Printf("code: %s", code) // 记录请求参数
	//4. 保存验证码到redis
	//5. 发送验证码
	successMsg := "success"      // 声明字符串变量
	resp.SetMessage(&successMsg) // 传递指针
	resp.SetSuccess(true)
	return resp, nil
}

func generateCaptcha() (string, error) {
	// 生成一个0到999999之间的随机数（6位数）
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // 1000000是10^6，即最大值+1
	if err != nil {
		return "", err
	}
	// 将随机数转换为字符串，并确保长度为6位（不足时前面补0）
	captcha := fmt.Sprintf("%06d", n) // %06d确保至少6位数字，不足时前面补0
	return captcha, nil
}
