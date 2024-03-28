package api

import (
	"crypto/md5"
	"fmt"
	"naive-admin/internal/inout"
	"naive-admin/internal/model"
	"naive-admin/pkg/db"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var User = &user{}

type user struct {
}

func (user) Add(c *gin.Context) {
	var req inout.AddUserReq
	if err := c.Bind(&req); err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	db.Dao.Transaction(func(tx *gorm.DB) error {
		var newUser = model.User{
			Username:   req.Username,
			Password:   fmt.Sprintf("%x", md5.Sum([]byte(req.Password))),
			Enable:     req.Enable,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}
		tx.Create(&model.Profile{
			UserId:   newUser.ID,
			NickName: newUser.Username,
		})
		for _, id := range req.RoleIds {
			tx.Create(&model.UserRolesRole{
				UserId: newUser.ID,
				RoleId: id,
			})
		}
		Resp.Succ(c, "")
		return nil
	})

}

func (user) Detail(c *gin.Context) {
	var res = &inout.UserDetailRes{}
	var uid, _ = c.Get("uid")
	db.Dao.Model(model.User{}).Where("id=?", uid).Find(&res)
	db.Dao.Model(model.Profile{}).Where("userId=?", uid).Find(&res.Profile)
	uroleIdList := db.Dao.Model(model.UserRolesRole{}).Where("userId=?", uid).Select("roleId")
	db.Dao.Model(model.Role{}).Where("id IN (?)", uroleIdList).Find(&res.Roles)
	if len(res.Roles) > 0 {
		res.CurrentRole = res.Roles[0]
	}
	Resp.Succ(c, res)

}
