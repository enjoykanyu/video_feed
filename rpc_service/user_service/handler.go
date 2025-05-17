package main

import (
	"context"
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
	return
}
