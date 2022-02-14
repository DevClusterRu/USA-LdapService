package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path"
)



type Config struct {
	Servers   map[string]struct{
		Login string `yaml:"login"`
		Password string `yaml:"password"`
		Urls []string `yaml:"urls"`
	} `yaml:"servers"`
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

	//config.LdapConnections = make(map[string]LdapCreds)
	//
	//for _, v:= range config.Servers {
	//	config.LdapConnections[v.Domain] = LdapCreds{
	//		v.Server,
	//		v.Login,
	//		v.Password,
	//		v.CN,
	//	}
	//}


	return config, err
}
