package internal

import (
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
		fmt.Fprintf(w, "Do user search")
	default:
		fmt.Fprintf(w, "Unknown command: %s", val)
	}

}

func NewLdapClient(server string) *LdapConnection {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	l, err := ldap.DialURL(server, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatal(err)
	}
	err = l.Bind(user, password)
	if err != nil {
		log.Fatal(err)
	}
	return &LdapConnection{l}
}

func (l *LdapConnection) LdapBind(user, password string) {
}

func (l *LdapConnection) LdapSearchUser(user, baseDN string) *ldap.SearchResult {
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(user))
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})
	result, err := l.Conn.Search(searchReq)
	if err != nil {
		fmt.Errorf("failed to query LDAP: %w", err)
	}
	return result
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
