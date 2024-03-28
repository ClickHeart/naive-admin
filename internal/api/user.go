package api

import (
	"naive-admin/internal/inout"
	"naive-admin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var User = &user{}

type user struct {
}

func (user) Register(c *gin.Context) {
	var req inout.AddUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	Resp.Succ(c, nil)
}

func (user) Add(c *gin.Context) {
	var req inout.AddUserReq
	if err := c.Bind(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := service.UserService.Register(c, req); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	Resp.Succ(c, nil)

}

func (user) Detail(c *gin.Context) {
	var uid, _ = c.Get("uid")
	userId, _ := uid.(int)
	data, err := service.UserService.Detail(c, userId)
	if err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	Resp.Succ(c, data)
}
