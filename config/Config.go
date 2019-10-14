package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	FileName string
}

func (c *Config) initConfig() error {
	if c.FileName != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.FileName)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
	}
	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")
	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// 监听配置文件是否改变,用于热更新
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed: %s\n", e.Name)
	})
}
func init() {
	c := Config{}
	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		log.Info(err)
	}
	c.watchConfig()
}

func GetValue(key string) string {
	return viper.GetString(key)
}

func GetValues(key string) []string {
	return viper.GetStringSlice(key)
}

func GetViperUnmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}
