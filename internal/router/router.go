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

	r.GET("/auth/captcha", api.Auth.Captcha)
	r.POST("/auth/register", api.Auth.Register)
	r.POST("/auth/login", api.Auth.Login)

	a := r.Group("/").Use(middleware.Jwt())
	{
		a.POST("/auth/logout", api.Auth.Logout)
		a.POST("/auth/password", api.Auth.Password)

		a.GET("/user", api.User.List)
		a.POST("/user", api.User.Add)
		a.DELETE("/user/:id", api.User.Delete)
		a.PATCH("/user/password/reset/:id", api.User.Update)
		a.PATCH("/user/:id", api.User.Update)
		a.PATCH("/user/profile/:id", api.User.Profile)
		a.GET("/user/detail", api.User.Detail)

		a.GET("/role", api.Role.List)
		a.POST("/role", api.Role.Add)
		a.PATCH("/role/:id", api.Role.Update)
		a.DELETE("/role/:id", api.Role.Delete)
		a.PATCH("/role/users/add/:id", api.Role.AddUser)
		a.PATCH("/role/users/remove/:id", api.Role.RemoveUser)
		a.GET("/role/page", api.Role.ListPage)
		a.GET("/role/permissions/tree", api.Role.PermissionsTree)
	}
}
