package main

import (
	"USALdapNewWave/internal"
	"log"
	"net/http"
)

func main() {

	l := internal.NewLdapClient("ldaps://DC-TEST.test.lab:636")
	l.LdapBind("test\\admin", "Te$t0vP@$$")

	http.HandleFunc("/ldap", l.LdapHandler)
	log.Println("Starting webserver...")
	err := http.ListenAndServe(":8085", nil)
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

	//result := l.SearchUser("23erwe", "OU=Тестовые пользователи,DC=test,DC=lab")
	//
	//log.Println("Got", len(result.Entries), "search results")
	//fmt.Println(result.Entries[0].DN)

}
