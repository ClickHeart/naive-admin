package service

import (
	"context"
	"errors"
	"naive-admin/internal/inout"
	"naive-admin/internal/model"
	r "naive-admin/internal/repository"
	"naive-admin/pkg/config"
	"naive-admin/pkg/utils/jwt"
	"naive-admin/pkg/utils/paginator"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var UserService = &userService{}

type userService struct {
}

func (userService) Create(c context.Context, data inout.AddUserReq) error {
	// check username
	user, err := r.UserRepo.GetByUsername(c, data.Username)
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

	err = r.Repo.Transaction(c, func(ctx context.Context) error {
		// Create a user
		if err = r.UserRepo.Create(ctx, &newUser); err != nil {
			return err
		}

		if err = r.ProfileRepo.Create(ctx, &model.Profile{
			UserId:   newUser.ID,
			NickName: newUser.Username,
		}); err != nil {
			return err
		}

		for _, id := range data.RoleIds {
			if err = r.UserRolesRepo.Create(ctx, &model.UserRolesRole{
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

func (userService) Login(c *gin.Context, data inout.LoginReq) (string, error) {
	user, err := r.UserRepo.GetByUsername(c, data.Username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return "", err
	}
	j := jwt.NewJwt(&config.Conf.Security.Jwt)
	token, err := j.GenToken(user.ID, time.Now().Add(time.Hour*1))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (userService) GetDetail(c context.Context, uid int) (data inout.UserDetailRes, err error) {

	user, err := r.UserRepo.GetById(c, uid)
	if err != nil {
		return
	}
	data.User = *user

	profile, err := r.ProfileRepo.GetByUserId(c, uid)
	if err != nil {
		return
	}
	data.Profile = profile

	roles, err := r.RoleRepo.GetByUserId(c, uid)
	if err != nil {
		return
	}
	data.Roles = *roles

	if len(data.Roles) > 0 {
		data.CurrentRole = data.Roles[0]
	}
	return
}

func (userService) ChangePassword(c context.Context, uid int, data inout.AuthPwReq) error {
	user, err := r.UserRepo.GetById(c, uid)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword)); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if err := r.UserRepo.Update(c, user); err != nil {
		return err
	}

	return nil
}

func (userService) ChangeProfile(c context.Context, data inout.PatchProfileUserReq) error {
	var profile = &model.Profile{
		ID:       data.Id,
		Gender:   data.Gender,
		Address:  data.Address,
		Email:    data.Email,
		NickName: data.NickName,
	}

	if err := r.ProfileRepo.Update(c, profile); err != nil {
		return err
	}

	return nil
}

func (userService) Update(c context.Context, data inout.PatchUserReq) error {
	var user model.User
	var profile model.Profile
	if data.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	if data.Enable != nil {
		user.Enable = *data.Enable
	}
	if data.Username != nil {
		user.Username = *data.Username
		profile.NickName = *data.Username
		profile.UserId = data.Id
		if err := r.ProfileRepo.Update(c, &profile); err != nil {
			return err
		}
	}
	if err := r.UserRepo.Update(c, &user); err != nil {
		return err
	}
	return nil
}

func (userService) Delete(c context.Context, uid int) (err error) {
	err = r.Repo.Transaction(c, func(ctx context.Context) error {
		// Create a user
		if err = r.UserRepo.DeleteById(ctx, uid); err != nil {
			return err
		}

		if err = r.UserRolesRepo.DeleteByUid(ctx, uid); err != nil {
			return err
		}

		if err = r.ProfileRepo.DeleteByUid(ctx, uid); err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (userService) GetList(c context.Context, query *map[string]string) (userList []inout.UserListItem, total int64, err error) {
	pageNo, _ := strconv.Atoi((*query)["pageNo"])
	pageSize, _ := strconv.Atoi((*query)["pageSize"])
	p := &paginator.Page[model.Profile]{
		CurrentPage: int64(pageNo),
		PageSize:    int64(pageSize),
	}
	if err := r.ProfileRepo.GetList(c, p, query); err != nil {
		return userList, total, err
	}
	total = p.Total
	for _, data := range p.Data {
		uInfo, err := r.UserRepo.GetById(c, data.UserId)
		if err != nil {
			return userList, total, err
		}
		rols, err := r.RoleRepo.GetByUserId(c, data.UserId)
		if err != nil {
			return userList, total, err
		}
		userList = append(userList, inout.UserListItem{
			ID:         uInfo.ID,
			Username:   uInfo.Username,
			Enable:     uInfo.Enable,
			CreateTime: uInfo.CreateTime,
			UpdateTime: uInfo.UpdateTime,
			Gender:     data.Gender,
			Avatar:     data.Avatar,
			Address:    data.Address,
			Email:      data.Email,
			Roles:      *rols,
		})
	}
	return
}
