package internal

import (
	"USALdapNewWave/config"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"log"
	"net/http"
)

type LdapConnection struct {
	Conn *ldap.Conn
}

func checkRequersStructure(need []string, got map[string]string) bool {
	for _, v := range need {
		if _, ok := got[v]; !ok {
			return false
		}
	}
	return true
}

func (l *LdapConnection) LdapHandler(w http.ResponseWriter, req *http.Request) {
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
	case "SearchUser":
		if !checkRequersStructure([]string{"filter", "baseDN"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		result := l.LdapSearchUser(params["filter"], params["baseDN"])
		b, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		fmt.Fprintf(w, string(b))
	case "CreateUser":
		fmt.Fprintf(w, "Do user search")
		if !checkRequersStructure([]string{"user", "password"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		result := l.LdapSearchUser(params["user"], params["baseDN"])
		fmt.Println(result)
	default:
		fmt.Fprintf(w, "Unknown command: %s", val)
	}

}

func NewLdapClient(c *config.Config) *LdapConnection {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	l, err := ldap.DialURL(c.Ldap_server, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatal(err)
	}
	err = l.Bind(c.Ldap_login, c.Ldap_password)
	if err != nil {
		log.Fatal(err)
	}
	return &LdapConnection{l}
}

func (l *LdapConnection) LdapSearchUser(filter, baseDN string) (res []string) {
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})
	result, err := l.Conn.Search(searchReq)
	if err != nil {
		fmt.Errorf("failed to query LDAP: %w", err)
		return
	}
	for _, v := range result.Entries {
		res = append(res, v.Attributes[0].Values[0])
	}
	return
}

func (l *LdapConnection) LdapCreateGroup(group, baseDN string) {
	addReq := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", group, baseDN), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "group"})
	addReq.Attribute("name", []string{"testgroup"})
	addReq.Attribute("sAMAccountName", []string{group})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	addReq.Attribute("groupType", []string{fmt.Sprintf("%d", 0x00000004|0x80000000)})

	if err := l.Conn.Add(addReq); err != nil {
		log.Fatal("error adding group:", addReq, err)
	}
	fmt.Println("DONE")
}

func (l *LdapConnection) LdapCreateNewUser(name, group string) {
	addReq := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", name, group), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "user", "person"})
	addReq.Attribute("name", []string{name})
	addReq.Attribute("sAMAccountName", []string{name})
	addReq.Attribute("userAccountControl", []string{fmt.Sprintf("%d", 0x0202)})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	addReq.Attribute("userPrincipalName", []string{fmt.Sprintf("%s@example.com", name)})
	addReq.Attribute("accountExpires", []string{fmt.Sprintf("%d", 0x00000000)})

	if err := l.Conn.Add(addReq); err != nil {
		log.Fatal("error adding service:", addReq, err)
	}
	fmt.Println("DONE")
}

func (l *LdapConnection) LdapChangeUserPassword(user, group, newpassword string) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, err := utf16.NewEncoder().String("\"" + newpassword + "\"")
	if err != nil {
		log.Fatal(err)
	}

	modReq := ldap.NewModifyRequest(fmt.Sprintf("CN=%s,%s", user, group), []ldap.Control{})
	modReq.Replace("unicodePwd", []string{pwdEncoded})
	if err := l.Conn.Modify(modReq); err != nil {
		log.Fatal("error setting user password:", modReq, err)
	}
	fmt.Println("DONE")
}

func (l *LdapConnection) LdapUserActivate(user, group string) {
	modReq := ldap.NewModifyRequest(fmt.Sprintf("CN=%s,%s", user, group), []ldap.Control{})
	modReq.Replace("userAccountControl", []string{fmt.Sprintf("%d", 0x0200)})

	if err := l.Conn.Modify(modReq); err != nil {
		log.Fatal("error enabling user account:", modReq, err)
	}
}

func (l *LdapConnection) LdapAssignUserToGroup(targetGroup, user string) {

	addReq := ldap.NewModifyRequest("CN=Администраторы домена,CN=Users,DC=test,DC=lab", []ldap.Control{})
	addReq.Add("member", []string{"CN=A.Kikos3,CN=Users,DC=test,DC=lab"})

	if err := l.Conn.Modify(addReq); err != nil {
		log.Fatal("error adding service:", addReq, err)
	}
	fmt.Println("DONE")

}
