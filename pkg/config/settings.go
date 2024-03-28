package config

import (
	"naive-admin/pkg/log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf *Config

func Init() (err error) {

	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}
	//设置读取的文件路径
	configPath := filepath.Join(path, "config")
	viper.AddConfigPath(configPath)
	//设置文件的类型
	viper.SetConfigType("yaml")

	// 读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Error(err)
	}
	// 监视配置信息
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			log.Error(err)
		}
		log.SetLevel(Conf.Log.Level)
	})
	return
}
