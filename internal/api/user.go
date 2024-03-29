package api

import (
	"naive-admin/internal/inout"
	s "naive-admin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var User = &user{}

type user struct {
}

func (user) Detail(c *gin.Context) {
	var uid, _ = c.Get("uid")
	userId, _ := uid.(int)
	data, err := s.UserService.GetDetail(c, userId)
	if err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	Resp.Succ(c, data)
}

func (user) Profile(c *gin.Context) {
	var req inout.PatchProfileUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.UserService.ChangeProfile(c, req); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	Resp.Succ(c, nil)
}

func (user) Update(c *gin.Context) {
	var req inout.PatchUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.UserService.Update(c, req); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	Resp.Succ(c, nil)
}
