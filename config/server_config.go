package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/hillguo/sanlog"
	"sync"
)

type Config struct {
	Server ServerConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

var config *Config
var config_mutex sync.Mutex

func GetConfig() *Config{
	return config
}

func ReloadServerConfig(){
	config_mutex.Lock()
	defer config_mutex.Unlock()

	var tmpConfig Config
	if _, err := toml.DecodeFile("server_config.toml", &config); err != nil {
		log.Error("reload server config failed", err.Error())
		return
	}
	config = &tmpConfig
}
func init(){
	config = &Config{}
	if _, err := toml.DecodeFile("server_config.toml", config); err != nil {
		fmt.Println(err)
		panic("config parse fail" + err.Error())
	}
	log.Info(config)
}