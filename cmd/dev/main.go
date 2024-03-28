package main

import (
	"naive-admin/pkg/config"
	"naive-admin/pkg/db"
	"naive-admin/pkg/log"
	"naive-admin/pkg/server"
)

func main() {

	// 加载配置
	if err := config.Init(); err != nil {
		log.Error("init settings failed")
		return
	}

	// 配置LOG
	log.SetLevel(config.Conf.Log.Level)

	// 初始化pgsql连接
	db.PgsqlInit()

	// 初始化gin
	app := server.Init()

	// 优雅退出
	server.Exit(app)

}
