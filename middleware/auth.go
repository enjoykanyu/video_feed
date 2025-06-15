package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
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

		// 验证Token有效性（需实现verifyToken函数）
		//if !verifyToken(token) {
		//	c.AbortWithStatusJSON(401, map[string]string{"error": "无效Token"})
		//	return
		//}
		c.Next(ctx)
	}
}
