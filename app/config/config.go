package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path"
)

type Config struct {
	Ldap_server   string `yaml:"ldap_server"`
	Ldap_login    string `yaml:"ldap_login"`
	Ldap_password string `yaml:"ldap_password"`
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

	return config, err
}
