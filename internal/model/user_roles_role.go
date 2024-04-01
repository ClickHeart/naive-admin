package model

type UserRoles struct {
	UserId int
	RoleId int
}

func (UserRoles) TableName() string {
	return "user_roles_role"
}
