package model

type UserRolesRole struct {
	UserId int
	RoleId int
}

func (UserRolesRole) TableName() string {
	return "user_roles_role"
}
