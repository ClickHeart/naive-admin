package repository

import (
	"context"
	"errors"
	"naive-admin/internal/model"
	"naive-admin/pkg/utils/paginator"

	"gorm.io/gorm"
)

var ProfileRepo = &profileRepo{}

type profileRepo struct {
	*Repository
}

func (r profileRepo) Create(c context.Context, profile *model.Profile) error {
	if err := r.DB(c).Create(profile).Error; err != nil {
		return err
	}
	return nil
}

func (r profileRepo) DeleteByUid(c context.Context, uid int) (err error) {
	if err := r.DB(c).Where("userId =?", uid).Delete(&model.Profile{}).Error; err != nil {
		return err
	}
	return
}

func (r profileRepo) Update(c context.Context, profile *model.Profile) error {

	if err := r.DB(c).Model(&profile).Updates(profile).Error; err != nil {
		return err
	}
	return nil
}

func (r profileRepo) GetByUserId(c context.Context, id int) (*model.Profile, error) {
	var profile model.Profile
	if err := r.DB(c).Model(model.Profile{}).Where("userId = ?", id).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r profileRepo) GetList(c context.Context, p *paginator.Page[model.Profile], query *map[string]string) (err error) {
	orm := r.DB(c).Model(model.Profile{})
	if (*query)["gender"] != "" {
		orm = orm.Where("gender=?", (*query)["gender"])
	}
	if (*query)["enable"] != "" {
		orm = orm.Where("userId in(?)", r.DB(c).Model(model.User{}).Where("enable=?", (*query)["enable"]).Select("id"))
	}
	if (*query)["username"] != "" {
		orm = orm.Where("nickName like ?", "%"+(*query)["username"]+"%")
	}
	if err := p.SelectPages(orm); err != nil {
		return err
	}
	return err
}
