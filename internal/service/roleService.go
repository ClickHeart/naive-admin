package service

import (
	"context"
	"naive-admin/internal/inout"
	"naive-admin/internal/model"
	r "naive-admin/internal/repository"
)

var RoleService = &roleService{}

type roleService struct {
}

func (roleService) GetList(c context.Context) (data *[]*model.Role, err error) {
	data, err = r.RoleRepo.GetList(c)
	return
}

func (roleService) Create(c context.Context, data inout.AddRoleReq) (err error) {

	err = r.Repo.Transaction(c, func(ctx context.Context) error {
		var record = model.Role{
			Code:   data.Code,
			Name:   data.Name,
			Enable: data.Enable,
		}
		if err := r.RoleRepo.Create(c, &record); err != nil {
			return err
		}

		for _, id := range data.PermissionIds {
			r.RolePermissionsrRepo.Create(c, &model.RolePermissions{
				RoleId:       record.ID,
				PermissionId: id,
			})
		}
		return nil
	})

	return

}
