package api

import (
	"naive-admin/internal/inout"
	s "naive-admin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Role = &role{}

type role struct {
}

func (role) List(c *gin.Context) {
	var res inout.RoleListRes
	data, err := s.RoleService.GetList(c)
	if err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
	}
	res = append(res, *data...)
	Resp.Succ(c, res)
}

func (role) Add(c *gin.Context) {
	var req inout.AddRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Resp.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.RoleService.Create(c, req); err != nil {
		Resp.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	Resp.Succ(c, nil)
}
