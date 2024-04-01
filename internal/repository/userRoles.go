package repository

import (
	"context"
	"naive-admin/internal/model"
)

var UserRolesRepo = &userRolesRepo{}

type userRolesRepo struct {
	*Repository
}

func (r userRolesRepo) Create(c context.Context, userRoles *model.UserRoles) error {
	if err := r.DB(c).Create(userRoles).Error; err != nil {
		return err
	}
	return nil
}

func (r userRolesRepo) DeleteByUid(c context.Context, uid int) (err error) {
	if err := r.DB(c).Where("userId =?", uid).Delete(&model.UserRoles{}).Error; err != nil {
		return err
	}
	return
}
