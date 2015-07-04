package config

import (
	"encoding/json"
	"io/ioutil"
	"github.com/toolkits/file"
	"sync"
	"log"
)

type GlobalConfig struct {
	HttpAddress string `json:"httpAddress"`
	Sender      string `json:"sender"`
	JobKey      string `json:"jobKey"`
	LogKey      string `json:"logKey"`
	Redis       RedisConfig `json:"redis"`
	Email       EmailConfig `json:"email"`
	Smtp        SmtpConfig `json:"smtp"`
	Sendcloud   SendcloudConfig `json:"sendcloud"`
}

type RedisConfig struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Timeout  string `json:"timeout"`
}

type SmtpConfig struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendcloudConfig struct {
	Address  string `json:"address"`
	Api_user string `json:"api_user"`
	Api_key  string `json:"api_key"`
}

type EmailConfig struct {
	From     string `json:"from"`
	Fromname string `json:"fromname"`
	Replyto  string `json:"replyto"`
	//Mailtype string `json:"mailtype"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// 初始化配置文件
func Parse(path string) {
	if !file.IsExist(path) {
		panic("Parse - parse file not exists")
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Parse - read file error ", err.Error())
	}

	ConfigFile = path

	var c GlobalConfig
	err = json.Unmarshal(buf, &c)
	if err != nil {
		log.Fatalln("Parse - Unmarshal ", ConfigFile, " failed, error :", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("Parse - parse config file success - ", ConfigFile)
}
