package service

import (
	"context"
	"errors"
	"naive-admin/internal/inout"
	"naive-admin/internal/model"
	"naive-admin/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var UserService = &userService{}

type userService struct {
}

func (userService) Register(c context.Context, data inout.AddUserReq) error {
	// check username
	user, err := repository.UserRepo.GetByUsername(c, data.Username)
	if err != nil {
		return err
	}

	if user != nil {
		return errors.New("username is already used")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var newUser = model.User{
		Username:   data.Username,
		Password:   string(hashedPassword),
		Enable:     data.Enable,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err = repository.Repo.Transaction(c, func(ctx context.Context) error {
		// Create a user
		if err = repository.UserRepo.Create(ctx, &newUser); err != nil {
			return err
		}

		if err = repository.ProfileRepo.Create(ctx, &model.Profile{
			UserId:   newUser.ID,
			NickName: newUser.Username,
		}); err != nil {
			return err
		}

		for _, id := range data.RoleIds {
			if err = repository.UserRolesRepo.Create(ctx, &model.UserRolesRole{
				UserId: newUser.ID,
				RoleId: id,
			}); err != nil {
				return err
			}
		}
		return nil
	})

	return nil

}

func (userService) Detail(c context.Context, uid int) (data inout.UserDetailRes, err error) {

	user, err := repository.UserRepo.GetById(c, uid)
	if err != nil {
		return
	}
	data.User = *user

	profile, err := repository.ProfileRepo.GetByUserId(c, uid)
	if err != nil {
		return
	}
	data.Profile = profile

	roles, err := repository.RoleRepo.GetByUserId(c, uid)
	if err != nil {
		return
	}
	data.Roles = *roles

	if len(data.Roles) > 0 {
		data.CurrentRole = data.Roles[0]
	}
	return
}
