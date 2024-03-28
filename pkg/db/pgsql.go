package db

import (
	"fmt"
	"naive-admin/pkg/config"
	"naive-admin/pkg/log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Dao *gorm.DB

var once sync.Once

func NewPgDao(cfg *config.Database) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DB,
		cfg.Port,
	)

	gormconfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	dao, err := gorm.Open(postgres.Open(dsn), gormconfig)
	if err != nil {
		log.Panic(err)
	}

	dbCon, err := dao.DB()
	if err != nil {
		log.Panic(err)
	}
	dbCon.SetMaxIdleConns(1)
	dbCon.SetMaxOpenConns(5)
	dbCon.SetConnMaxLifetime(time.Hour)
	log.Info("Pgsql is Connected")

	return dao

}

func PgsqlInit() {

	singleton := func() {
		if config.Conf == nil {
			log.Panic("config is nil!")
		}
		Dao = NewPgDao(config.Conf.Data.Pgsql)
	}
	once.Do(singleton)
}
