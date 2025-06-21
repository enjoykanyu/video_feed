package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

const (
	tokenPrefix     = "user_token:"
	tokenExpiration = 24 * time.Hour
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 放行登录/注册接口
		if strings.HasPrefix(c.FullPath(), "/douyin/user/register") ||
			strings.HasPrefix(c.FullPath(), "/douyin/user/login") {
			c.Next(ctx)
			return
		}

		// 检查Token字段
		token := c.GetHeader("Authorization")
		if len(token) == 0 {
			c.AbortWithStatusJSON(401, map[string]string{"error": "未授权访问"})
			return
		}
		ttl, err := redis.Client{}.TTL(ctx, tokenPrefix+string(token)).Result()
		if err != nil {
			return
		}

		// 已过期直接返回错误
		if ttl < 0 {
			return
		}
		err = redis.Client{}.Expire(ctx, string(token), tokenExpiration).Err()

		if err != nil {
			c.AbortWithStatusJSON(401, map[string]string{"error": "无效Token"})
			return
		}
		c.Next(ctx)
	}
}
