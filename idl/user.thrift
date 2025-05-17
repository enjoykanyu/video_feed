namespace go user

// 用户注册请求
struct RegisterRequest {
    1: required string phone;       // 手机号
    2: required string verifyCode;  // 验证码
}

// 用户注册响应
struct RegisterResponse {
    1: required bool success;       // 是否成功
    2: optional string message;     // 错误信息
    3: optional i64 userId;         // 用户ID
    4: optional string username;    // 自动生成的用户名
}

// 用户登录请求
struct LoginRequest {
    1: required string phone;       // 手机号
    2: required string verifyCode;  // 验证码
}

// 用户登录响应
struct LoginResponse {
    1: required bool success;       // 是否成功
    2: optional string message;     // 错误信息
    3: optional i64 userId;         // 用户ID
    4: optional string token;       // 登录令牌
    5: optional UserInfo userInfo;  // 用户信息
}

// 用户基本信息
struct UserInfo {
    1: required i64 userId;         // 用户ID
    2: required string username;    // 用户名
    3: optional string avatar;      // 头像URL
    4: optional string signature;   // 个性签名
    5: optional i32 followCount;    // 关注数
    6: optional i32 followerCount;  // 粉丝数
    7: optional bool isFollow;      // 当前用户是否关注了该用户
}

// 获取用户信息请求
struct GetUserInfoRequest {
    1: required i64 userId;         // 要查询的用户ID
    2: optional i64 currentUserId;  // 当前登录的用户ID
}

// 获取用户信息响应
struct GetUserInfoResponse {
    1: required bool success;       // 是否成功
    2: optional string message;     // 错误信息
    3: optional UserInfo userInfo;  // 用户信息
}

// 发送验证码请求
struct SendVerifyCodeRequest {
    1: required string phone;       // 手机号
    2: required i32 codeType;       // 验证码类型: 1-注册，2-登录，3-重置密码
}

// 发送验证码响应
struct SendVerifyCodeResponse {
    1: required bool success;       // 是否成功
    2: optional string message;     // 错误信息
}

// 用户服务接口定义
service UserService {
    // 用户注册
    RegisterResponse Register(1: RegisterRequest req);
    
    // 用户登录
    LoginResponse Login(1: LoginRequest req);
    
    // 获取用户信息
    GetUserInfoResponse GetUserInfo(1: GetUserInfoRequest req);
    
    // 发送验证码
    SendVerifyCodeResponse SendVerifyCode(1: SendVerifyCodeRequest req);
} 