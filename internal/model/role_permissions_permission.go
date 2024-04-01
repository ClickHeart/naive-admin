package model

type RolePermissions struct {
	RoleId       int `gorm:"column:roleId"`
	PermissionId int `gorm:"column:permissionId"`
}

func (RolePermissions) TableName() string {
	return "role_permissions_permission"
}
