package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"
	"video_douyin/dal/model"
	user "video_douyin/kitex_gen/user"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

const (
	verifyCodePrefix = "verify_code:"
	tokenPrefix      = "user_token:"
	codeExpiration   = 5 * time.Minute
	tokenExpiration  = 24 * time.Hour
)

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// TODO: Your code here...

	resp = user.NewRegisterResponse()
	if resp == nil {
		return nil, errors.New("failed to initialize response")
	}
	// 1. Validate the request parameters
	log.Printf("进入rpc: %v\n", req.GetPhone())
	if req.GetPhone() == "" || req.GetVerifyCode() == "" {
		msg := "参数为空错误" // 改为包级变量或堆分配
		resp.SetMessage(&msg)
		resp.SetSuccess(false)
		if resp.Message != nil {
			log.Printf("错误详情: %s\n", *resp.Message)
		}
		return
	}

	// 2. 验证码校验
	storedCode, err := s.redis.Get(ctx, verifyCodePrefix+req.GetPhone()).Result()
	if err != nil || storedCode != req.GetVerifyCode() {
		msg := "验证码错误或已过期"
		resp.SetMessage(&msg)
		resp.SetSuccess(false)
		if resp.Message != nil {
			log.Printf("错误详情: %s\n", *resp.Message)
		}
		return resp, nil
	}

	// 3. 检查用户是否存在 存在则进入注册网页
	var existingUser model.User
	if err := s.db.Where("phone = ?", req.GetPhone()).First(&existingUser).Error; err == nil {
		errorMsg := "用户已存在"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 4. 创建用户
	newUser := model.User{
		Phone:    req.GetPhone(),
		Nickname: fmt.Sprintf("用户%s", req.GetPhone()[:4]),
	}
	if err := s.db.Create(&newUser).Error; err != nil {
		errorMsg := "注册失败"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 5. 生成token
	token, err := generateToken()
	if err != nil {
		errorMsg := "系统错误"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 6. 存储token
	if err := s.redis.Set(tokenPrefix+token, newUser.ID, tokenExpiration).Err(); err != nil {
		errorMsg := "系统错误"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	successMsg := "注册成功"
	resp.SetMessage(&successMsg)
	resp.SetSuccess(true)
	// resp.SetToken(token)// 需改下idl
	return resp, nil
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
	//2，验证码校验
	// 2. 验证码校验
	storedCode, err := s.redis.Get(verifyCodePrefix + req.GetPhone()).Result()
	if err != nil || storedCode != req.GetVerifyCode() {
		errorMsg := "验证码错误或已过期"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}
	// 3. 查询用户是否存在
	var existingUser model.User
	if err := s.db.Where("phone = ?", req.GetPhone()).First(&existingUser).Error; err != nil {
		errorMsg := "用户不存在"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 4. 生成token
	token, err := generateToken()
	if err != nil {
		errorMsg := "系统错误"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 5. 存储token
	if err := s.redis.Set(tokenPrefix+token, existingUser.ID, tokenExpiration).Err(); err != nil {
		errorMsg := "系统错误"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}
	successMsg := "登录成功"
	resp.SetMessage(&successMsg)
	resp.SetSuccess(true)
	// resp.SetToken(token) 需加上返回给前端
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.

func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	resp = user.NewGetUserInfoResponse()

	// 1. 参数校验
	if req.GetUserId() == 0 {
		errorMsg := "Invalid user id"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 2. 查询用户信息
	var userModel model.User
	if err := s.db.Where("id = ?", req.GetUserId()).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := "User not found"
			resp.SetMessage(&errorMsg)
			resp.SetSuccess(false)
			return resp, nil
		}
		errorMsg := "Database error"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}

	// 3. 构造UserInfo结构体
	userInfo := &user.UserInfo{
		UserId:   userModel.ID,
		Username: userModel.Nickname,
		Avatar:   &userModel.Avatar,
	}

	// 4. 设置响应
	resp.SetSuccess(true)
	resp.SetUserInfo(userInfo)
	return resp, nil
}

// SendVerifyCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) SendVerifyCode(ctx context.Context, req *user.SendVerifyCodeRequest) (resp *user.SendVerifyCodeResponse, err error) {
	// TODO: Your code here...
	resp = user.NewSendVerifyCodeResponse()
	phone := req.GetPhone()
	// 1. 校验参数异常
	if phone == "" {

		errorMsg := "Invalid request parameters" // 声明字符串变量
		resp.SetMessage(&errorMsg)               // 传递指针
		resp.SetSuccess(false)
	}
	//2,查询db看用户手机号是否存在
	CheckPhoneExists(s.db, phone)
	//3, 生成验证码
	code, _ := generateCaptcha()

	log.Printf("code: %s", code) // 记录请求参数
	//4. 保存验证码到redis
	if err := s.saveCodeToRedis(phone, code); err != nil {
		errorMsg := "Failed to save verification code"
		resp.SetMessage(&errorMsg)
		resp.SetSuccess(false)
		return resp, nil
	}
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

// 生成token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// CheckPhoneExists 检查手机号是否已存在
func CheckPhoneExists(db *gorm.DB, phone string) (bool, error) {
	var user = &model.User{}
	result := db.Where("phone = ?", phone).First(user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil // 手机号不存在
		}
		return false, result.Error // 查询出错
	}
	return true, nil // 手机号已存在
}

// 存入redis
func (s *UserServiceImpl) saveCodeToRedis(phone, code string) error {
	key := verifyCodePrefix + phone
	return s.redis.Set(key, code, codeExpiration).Err()
}
