package main

import (
	"USALdapNewWave/config"
	"fmt"
	"github.com/kardianos/osext"
	"log"
	"net/http"
	"os"
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

	//fmt.Println(config.Cfg)

	config.Cfg.LdapChangeUserPassword("cl.local","CN=Лесовский тест,OU=AХВ,OU=Clients,DC=cl,DC=local","Qwe12345678!@#")


	fmt.Scanln()
	os.Exit(0)


	http.HandleFunc("/ldap", config.Cfg.LdapHandler)
	log.Println("Starting webserver...")
	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}

	//10.200.88.11 test\admin:Te$t0vP@$$

	//l := internal.NewLdapClient("ldaps://DC-TEST.test.lab:636")
	//l.LdapBind("test\\admin", "Te$t0vP@$$")

	//l.LdapCreateNewUser("A.Kikos3", "CN=Users,DC=test,DC=lab")
	//l.LdapChangeUserPassword("A.Kikos3", "CN=Users,DC=test,DC=lab", "xX1234567")
	//l.LdapUserActivate("A.Kikos3", "CN=Users,DC=test,DC=lab")

	//l.LdapAssignUserToGroup()

	//l.LdapReChangeUserPassword("A.Kikos3", "CN=Users,DC=test,DC=lab")

	//l.LdapCreateGroup("anotherGroup", "OU=Тестовые пользователи,DC=test,DC=lab")

	//result := l.LdapSearchUser("(objectClass=user)", "OU=Тестовые пользователи,DC=test,DC=lab")

	//fmt.Println(result)

	//
	//log.Println("Got", len(result.Entries), "search results")
	//fmt.Println(result.Entries[0].DN)

}
