package api

import (
	"naive-admin/internal/inout"
	s "naive-admin/internal/service"
	"naive-admin/pkg/utils/captcha"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var Auth = &auth{}

type auth struct {
}

func (auth) Captcha(c *gin.Context) {
	svg, code := captcha.GenerateSVG(80, 40)
	session := sessions.Default(c)
	session.Set("captch", code)
	session.Save()
	// 设置 Content-Type 为 "image/svg+xml"
	c.Header("Content-Type", "image/svg+xml; charset=utf-8")
	// 返回验证码
	c.Data(http.StatusOK, "image/svg+xml", svg)
}

func (auth) Register(c *gin.Context) {
	var req inout.AddUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.UserService.Create(c, req); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	Resp.Succ(c, nil)
}

func (auth) Login(c *gin.Context) {
	var req inout.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	session := sessions.Default(c)
	if req.Captcha != session.Get("captch") {
		Resp.Err(c, 20001, "验证码不正确")
		return
	}

	token, err := s.UserService.Login(c, req)
	if err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	Resp.Succ(c, inout.LoginRes{
		AccessToken: token,
	})
}

func (auth) Logout(c *gin.Context) {
	Resp.Succ(c, true)
}

func (auth) Password(c *gin.Context) {
	var req inout.AuthPwReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	uid, _ := c.Get("uid")
	userId, _ := uid.(int)
	if err := s.UserService.ChangePassword(c, userId, req); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	Resp.Succ(c, true)
}
