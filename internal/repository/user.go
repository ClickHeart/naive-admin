package repository

import (
	"context"
	"errors"
	"naive-admin/internal/model"

	"gorm.io/gorm"
)

var UserRepo = &userRepo{}

type userRepo struct {
	*Repository
}

func (r userRepo) Create(c context.Context, user *model.User) error {
	if err := r.DB(c).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepo) DeleteById(c context.Context, uid int) (err error) {
	if err := r.DB(c).Where("id =?", uid).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return err
}

func (r userRepo) Update(c context.Context, user *model.User) error {

	if err := r.DB(c).Model(&user).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepo) GetById(c context.Context, id int) (*model.User, error) {
	var user model.User
	if err := r.DB(c).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r userRepo) GetByUsername(c context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.DB(c).Model(model.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
