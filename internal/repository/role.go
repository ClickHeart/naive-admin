package repository

import (
	"context"
	"errors"
	"naive-admin/internal/model"

	"gorm.io/gorm"
)

var RoleRepo = &roleRepo{}

type roleRepo struct {
	*Repository
}

func (r roleRepo) Create(c context.Context, role *model.Role) error {
	if err := r.DB(c).Create(role).Error; err != nil {
		return err
	}
	return nil
}

func (r roleRepo) GetByUserId(c context.Context, uid int) (*[]*model.Role, error) {
	var roles []*model.Role
	roleIdList := r.DB(c).Model(model.UserRoles{}).Where("userId=?", uid).Select("roleId")
	if err := r.DB(c).Model(model.Role{}).Where("id IN (?)", roleIdList).Find(&roles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &roles, nil
}

func (r roleRepo) GetList(c context.Context) (data *[]*model.Role, err error) {
	if err := r.DB(c).Model(model.Role{}).Find(data).Error; err != nil {
		return data, err
	}
	return
}
