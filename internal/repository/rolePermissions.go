package repository

import (
	"context"
	"naive-admin/internal/model"
)

var RolePermissionsrRepo = &rolePermissonsRepo{}

type rolePermissonsRepo struct {
	*Repository
}

func (r rolePermissonsRepo) Create(c context.Context, rolePermissons *model.RolePermissions) (err error) {
	if err := r.DB(c).Create(rolePermissons).Error; err != nil {
		return err
	}
	return nil
}

func (r rolePermissonsRepo) DeleteByRoleId(c context.Context, id int) (err error) {
	if err := r.DB(c).Where("roleId =?", id).Delete(&model.RolePermissions{}).Error; err != nil {
		return err
	}
	return
}
