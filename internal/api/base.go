package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Resp = &rps{}

type rps struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
	OriginUrl string `json:"originUrl"`
}

func (rps) Succ(c *gin.Context, data any) {
	resp := rps{
		Code:      0,
		Message:   "OK",
		Data:      data,
		OriginUrl: c.Request.URL.Path,
	}
	c.Set("succ_response", resp)
	c.JSON(http.StatusOK, resp)
}

func (rps) Err(c *gin.Context, ErrCode int, messge string) {
	resp := rps{
		Code:      ErrCode,
		Error:     "Error",
		Message:   messge,
		OriginUrl: c.Request.URL.Path,
	}
	c.Set("err_response", resp)
	c.JSON(http.StatusOK, resp)
}
