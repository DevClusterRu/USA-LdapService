package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sync"
)

type Config struct {
	CALLER_DC_WEBHOOK_PROD    string         `yaml:"CALLER_DC_WEBHOOK_PROD"`
	CALLER_DC_WEBHOOK_EXAMPLE string         `yaml:"CALLER_DC_WEBHOOK_EXAMPLE"`
	CALLER_DC_WEBHOOK         string         `yaml:"CALLER_DC_WEBHOOK"`
	DB_MONGO                  string         `yaml:"DB_MONGO"`
	MONGO_CONNECTION_PP       string         `yaml:"MONGO_CONNECTION_PP"`
	MONGO_CONNECTION_PROD     string         `yaml:"MONGO_CONNECTION_PROD"`
	DB_HOST                   string         `yaml:"DB_HOST"`
	DB_PORT                   int            `yaml:"DB_PORT"`
	DB_DATABASE               string         `yaml:"DB_DATABASE"`
	DB_USERNAME               string         `yaml:"DB_USERNAME"`
	DB_PASSWORD               string         `yaml:"DB_PASSWORD"`
	APIKEY                    string         `yaml:"APIKEY"`
	APP_URL                   string         `yaml:"APP_URL"`
	AWS_ACCESS_KEY_ID         string         `yaml:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY     string         `yaml:"AWS_SECRET_ACCESS_KEY"`
	AWS_DEFAULT_REGION        string         `yaml:"AWS_DEFAULT_REGION"`
	AWS_BUCKET                string         `yaml:"AWS_BUCKET"`
	AWS_URL                   string         `yaml:"AWS_URL"`
	TARANTOOL                 string         `yaml:"TARANTOOL"`
	TUSER                     string         `yaml:"TUSER"`
	TPASSWORD                 string         `yaml:"TPASSWORD"`
	ELASTIC_AUTH              string         `yaml:"ELASTIC_AUTH"`
	ELASTIC_POINT             string         `yaml:"ELASTIC_POINT"`
	Debug                     string         `yaml:"debug"`
	MID_UID                   map[string]int //messageId_100[5] 5 - count of hooks
}

var Cfg *Config

func NewConfig(fileName string) (config *Config, err error) {
	log.Printf("reading config from '%s'", fileName)
	if ext := path.Ext(fileName); ext != ".yaml" && ext != ".yml" {
		err = fmt.Errorf("invalid file '%s' extenstion, expected 'yaml' or 'yml'", ext)
		return
	}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		err = fmt.Errorf("can't read file '%s'", fileName)
		return
	}

	config = new(Config)
	if err = yaml.Unmarshal(file, config); err != nil {
		err = fmt.Errorf("file %s yaml unmarshal error: %v", fileName, err)
	}
	fmt.Println("Start generate DC key...")
	telnyxBot.SmsBot = new(telnyxBot.SmsBotStruct)
	if err = yaml.Unmarshal(file, telnyxBot.SmsBot); err != nil {
		err = fmt.Errorf("file %s yaml unmarshal error: %v", fileName, err)
	}
	telnyxBot.SmsBot.MessageIdAssigned = make(map[string]telnyxBot.MessageIdPack)
	telnyxBot.SmsBot.UserPhoneAnswersCounter = make(map[string]int)
	telnyxBot.SmsBot.Mutex = sync.Mutex{}

	go telnyxBot.SmsBot.PostDCKeyGen()
	fmt.Println("Done")

	config.MongoPP = &databases.MongoStructure{}
	config.MongoPP.Client = databases.NewMongoClient(config.MONGO_CONNECTION_PP)
	config.MongoPP.DB = config.DB_MONGO

	config.MongoPROD = &databases.MongoStructure{}
	config.MongoPROD.Client = databases.NewMongoClient(config.MONGO_CONNECTION_PROD)
	config.MongoPROD.DB = config.DB_MONGO

	config.MID_UID = make(map[string]int)

	config.Tarantool = databases.NewTarantool(config.TARANTOOL, config.TUSER, config.TPASSWORD)

	config.Mysql = databases.NewMysqlClient(fmt.Sprintf("%v:%v@(%v)/%v", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE))

	config.WLimitMutex = &sync.Mutex{}

	return config, err
}

//func (c *Config) HardDevicesInfo()  {
//	for {
//		for _, v := range c.Clouds {
//			token := dcGetDevicesStatuses.GetTokenAccess(v.Url_access, v.Login, v.Password)
//			devResponse := dcGetDevicesStatuses.GetDevicesInfo(v.Url_devices, token)
//			if devResponse.Result.Models != nil {
//				for _, v := range devResponse.Result.Models {
//					if !v.Available {
//						c.Metrics.Mutex.Lock()
//						c.Metrics.DeviceInfo[v.Id] = 1
//						c.Metrics.Mutex.Unlock()
//					}
//				}
//			}
//		}
//		time.Sleep(30*time.Minute)
//	}
//}
