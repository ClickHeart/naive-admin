package middleware

import (
	"naive-admin/internal/api"
	"naive-admin/pkg/config"
	"naive-admin/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			api.Resp.Err(c, 10002, "请求未携带token，无权限访问")
			c.Abort()
			return
		}

		j := jwt.NewJwt(&config.Conf.Security.Jwt)
		claims, err := j.ParseToken(token)
		if err != nil {
			api.Resp.Err(c, 10002, err.Error())
			c.Abort()
			return
		}

		c.Set("uid", claims.UserId)
		c.Next()
	}
}
