package router

import (
	"naive-admin/internal/api"
	"naive-admin/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 使用 cookie 存储会话数据
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("captch"))))
	r.Use(middleware.Cors())

	noAuth := r.Group("/")
	{
		noAuth.GET("/auth/captcha", api.Auth.Captcha)
		noAuth.POST("/auth/register", api.Auth.Register)
		noAuth.POST("/auth/login", api.Auth.Login)
	}

	Auth := r.Group("/").Use(middleware.Jwt())
	{
		Auth.POST("/auth/logout", api.Auth.Logout)
		Auth.POST("/auth/password", api.Auth.Password)

		Auth.GET("/user", api.User.List)
		Auth.POST("/user", api.User.Add)
		Auth.DELETE("/user/:id", api.User.Delete)
		Auth.PATCH("/user/password/reset/:id", api.User.Update)
		Auth.PATCH("/user/:id", api.User.Update)
		Auth.PATCH("/user/profile/:id", api.User.Profile)
		Auth.GET("/user/detail", api.User.Detail)
	}
}
