package main

import (
	"USA-LdapService/config"
	"github.com/kardianos/osext"
	"log"
	"net/http"
)

func main() {

	cfgfile := "local.yml"
	var err error

	config.Cfg, err = config.NewConfig(cfgfile)
	if err != nil {
		log.Println("Something wrong with config: ", err)
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			log.Fatal("Cant find env folder path", err)
		}
		config.Cfg, err = config.NewConfig(folderPath)
		if err != nil {
			log.Fatal("Config create error", err)
		}
	}

	http.HandleFunc("/ldap", config.Cfg.LdapHandler)
	log.Println("Starting webserver...")
	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}


}
