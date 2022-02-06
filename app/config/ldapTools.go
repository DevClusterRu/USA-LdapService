package config

import (
	"USALdapNewWave/randomHash"
	"encoding/json"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"log"
	"net/http"
	"time"
)




func checkRequersStructure(need []string, got map[string]string) bool {
	for _, v := range need {
		if _, ok := got[v]; !ok {
			return false
		}
	}
	return true
}

func (c *Config) LdapHandler(w http.ResponseWriter, req *http.Request) {
	var params map[string]string
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var val string
	var ok bool
	if val, ok = params["command"]; !ok {
		http.Error(w, "Attribute 'command' expected", 400)
		return
	}

	switch val {
	//case "SearchUser":
	//	if !checkRequersStructure([]string{"filter", "baseDN"}, params) {
	//		fmt.Fprintf(w, "Wrong parameters, good bye", val)
	//		return
	//	}
	//	result := l.LdapSearchUser(params["filter"], params["baseDN"])
	//	b, err := json.Marshal(result)
	//	if err != nil {
	//		fmt.Fprintf(w, err.Error())
	//	}
	//	fmt.Fprintf(w, string(b))
	//case "CreateUser":
	//	fmt.Fprintf(w, "Do user search")
	//	if !checkRequersStructure([]string{"user", "password"}, params) {
	//		fmt.Fprintf(w, "Wrong parameters, good bye", val)
	//		return
	//	}
	//	result := l.LdapSearchUser(params["user"], params["baseDN"])
	//	fmt.Println(result)
	case "DropPassword":
		fmt.Fprintf(w, "Do drop password: ")
		if !checkRequersStructure([]string{"user", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}

		newPassword := ""
		i := 0
		for i = 1; i < 4; i++ {
			newPassword = randomHash.RandomString(10)
			fmt.Println("Attempt #", i)
			if c.LdapChangeUserPassword(params["domain"], params["user"], newPassword) == true {
				break
			}
			fmt.Println("==>", newPassword)
			time.Sleep(100 * time.Millisecond)
		}
		if i >= 5 {
			fmt.Fprintf(w, "Max attempts! Exiting...")
			return
		}
		fmt.Fprintf(w, newPassword)
	default:
		fmt.Fprintf(w, "Unknown command: %s", val)
	}

}

func (c *Config) GetConn(server string) (*ldap.Conn, error) {
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true,
	//}

	//cert, err := tls.LoadX509KeyPair("PublicKey.crt","PrivateKey.key")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tlsCong := &tls.Config{Certificates: []tls.Certificate{cert}}

	//tlsCong := &tls.Config{
	//	InsecureSkipVerify: true,
	//}
	//tlsConfig := &tls.Config{InsecureSkipVerify: true}
	conn, err := ldap.DialURL(c.Servers[server].Urls[0])
	log.Println(c.Servers[server].Urls[0])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = conn.Bind(c.Servers[server].Login, c.Servers[server].Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conn, err
}

//func LdapSearchUser(filter, baseDN string) (res []string) {
//	//conn := l.GetConn()
//	//if conn == nil {
//	//	fmt.Println("Error when bind!")
//	//	return
//	//}
//	//searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})
//	//result, err := conn.Search(searchReq)
//	//if err != nil {
//	//	fmt.Errorf("failed to query LDAP: %w", err)
//	//	return
//	//}
//	//for _, v := range result.Entries {
//	//	res = append(res, v.Attributes[0].Values[0])
//	//}
//	//return
//}

func LdapCreateGroup(group, baseDN string) {
	//addReq := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", group, baseDN), []ldap.Control{})
	//addReq.Attribute("objectClass", []string{"top", "group"})
	//addReq.Attribute("name", []string{"testgroup"})
	//addReq.Attribute("sAMAccountName", []string{group})
	//addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	//addReq.Attribute("groupType", []string{fmt.Sprintf("%d", 0x00000004|0x80000000)})
	//
	//if err := l.Conn.Add(addReq); err != nil {
	//	log.Fatal("error adding group:", addReq, err)
	//}
	//fmt.Println("DONE")
}

func LdapCreateNewUser(name, group string) {
	//addReq := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", name, group), []ldap.Control{})
	//addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "user", "person"})
	//addReq.Attribute("name", []string{name})
	//addReq.Attribute("sAMAccountName", []string{name})
	//addReq.Attribute("userAccountControl", []string{fmt.Sprintf("%d", 0x0202)})
	//addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	//addReq.Attribute("userPrincipalName", []string{fmt.Sprintf("%s@example.com", name)})
	//addReq.Attribute("accountExpires", []string{fmt.Sprintf("%d", 0x00000000)})
	//
	//if err := l.Conn.Add(addReq); err != nil {
	//	log.Fatal("error adding service:", addReq, err)
	//}
	//fmt.Println("DONE")
}

func (c *Config) LdapChangeUserPassword(domain, user, newpassword string) bool {

	conn, err:= c.GetConn(domain)
	if err!=nil{
		log.Println(err)
		return false
	}
	defer conn.Close()
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

	pwdEncoded, err := utf16.NewEncoder().String("\"" + newpassword + "\"")
	if err != nil {
		log.Println(err)
		return false
	}

	modReq := ldap.NewModifyRequest(user, []ldap.Control{})
	modReq.Replace("unicodePwd", []string{pwdEncoded})
	if err := conn.Modify(modReq); err != nil {
		log.Println("error setting user password:", modReq, err)
		return false
	}

	fmt.Println("DONE")
	return true
}

//func (l *LdapConnection) LdapUserActivate(user, group string) {
//	modReq := ldap.NewModifyRequest(fmt.Sprintf("CN=%s,%s", user, group), []ldap.Control{})
//	modReq.Replace("userAccountControl", []string{fmt.Sprintf("%d", 0x0200)})
//
//	if err := l.Conn.Modify(modReq); err != nil {
//		log.Fatal("error enabling user account:", modReq, err)
//	}
//}
//
//func (l *LdapConnection) LdapAssignUserToGroup(targetGroup, user string) {
//
//	addReq := ldap.NewModifyRequest("CN=Администраторы домена,CN=Users,DC=test,DC=lab", []ldap.Control{})
//	addReq.Add("member", []string{"CN=A.Kikos3,CN=Users,DC=test,DC=lab"})
//
//	if err := l.Conn.Modify(addReq); err != nil {
//		log.Fatal("error adding service:", addReq, err)
//	}
//	fmt.Println("DONE")
//
//}
