package api

import (
	"naive-admin/internal/inout"
	s "naive-admin/internal/service"
	"net/http"
	"strconv"

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

func (user) Add(c *gin.Context) {
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

func (user) Delete(c *gin.Context) {
	uid := c.Param("id")
	userId, err := strconv.Atoi(uid)
	if err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
	}
	if err := s.UserService.Delete(c, userId); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	Resp.Succ(c, nil)
}

func (user) List(c *gin.Context) {
	var res inout.UserListRes
	query := map[string]string{
		"gender":   c.DefaultQuery("gender", ""),
		"enable":   c.DefaultQuery("enable", ""),
		"username": c.DefaultQuery("username", ""),
		"pageNo":   c.DefaultQuery("pageNo", "1"),
		"pageSize": c.DefaultQuery("pageSize", "10"),
	}

	data, total, err := s.UserService.GetList(c, &query)
	if err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	res.PageData = data
	res.Total = total
	Resp.Succ(c, res)
}
