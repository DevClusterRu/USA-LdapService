package config

import (
	"USA-LdapService/randomHash"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"log"
	"net/http"
	"time"
)

type Answer struct {
	Status   bool   `json:"result"`
	Response string `json:"response"`
}

func checkRequersStructure(need []string, got map[string]string) bool {
	for _, v := range need {
		if _, ok := got[v]; !ok {
			return false
		}
	}
	return true
}

func answerPack(status bool, response string) string {
	a := Answer{
		Status:   status,
		Response: response,
	}
	b, _ := json.Marshal(a)
	return string(b)
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
		if !checkRequersStructure([]string{"user", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		newPassword := ""
		errors := ""
		i := 0
		for i = 1; i < 4; i++ {
			newPassword = randomHash.RandomString(10)
			fmt.Println("Attempt #", i)
			err := c.LdapChangeUserPassword(params["domain"], params["user"], newPassword)
			if err == nil {
				fmt.Fprintf(w, answerPack(true, newPassword))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "CreateUser":
		if !checkRequersStructure([]string{"name", "baseDN", "domain", "phone", "email"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapCreateNewUser(params["name"], params["phone"], params["email"], params["baseDN"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "CreateOrganization":
		if !checkRequersStructure([]string{"group", "baseDN", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapCreateOrganization(params["group"], params["baseDN"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "CreateGroup":
		if !checkRequersStructure([]string{"group", "baseDN", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapCreateGroup(params["group"], params["baseDN"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "AssignUser":
		if !checkRequersStructure([]string{"targetGroup", "user", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapAssignUserToGroup(params["targetGroup"], params["user"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "UnassignUser":
		if !checkRequersStructure([]string{"targetGroup", "user", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors:=""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err:= c.LdapUnAssignUserToGroup(params["targetGroup"],params["user"],params["domain"])
			if err==nil{
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors+="\n"+err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "DeleteObject":
		if !checkRequersStructure([]string{"name", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapDeleteObject(params["name"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "UserActivation":
		if !checkRequersStructure([]string{"name", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapUserActivate(params["name"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "UserDisabling":
		if !checkRequersStructure([]string{"name", "domain"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors := ""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err := c.LdapUserDeactivate(params["name"], params["domain"])
			if err == nil {
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors += "\n" + err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))
	case "EditUser":
		if !checkRequersStructure([]string{"name", "baseDN", "domain", "phone", "email"}, params) {
			fmt.Fprintf(w, "Wrong parameters, good bye", val)
			return
		}
		i := 0
		errors:=""
		for i = 1; i < 4; i++ {
			fmt.Println("Attempt #", i)
			err:= c.LdapEditUser(params["name"],params["phone"],params["email"],params["baseDN"],params["domain"])
			if err==nil{
				fmt.Fprintf(w, answerPack(true, "success"))
				return
			}
			errors+="\n"+err.Error()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintf(w, answerPack(false, errors))

	default:
		fmt.Fprintf(w, "Unknown command: %s", val)
	}

}

func (c *Config) GetConn(server string) (*ldap.Conn, error) {
	if _, ok := c.Servers[server]; !ok {
		return nil, fmt.Errorf("Wrong server name: ", server)
	}
	conn, err := ldap.DialURL(c.Servers[server].Urls[0], ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
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

func (c *Config) LdapSendAddRequest(addReq *ldap.AddRequest, domain string) error {
	conn, err := c.GetConn(domain)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Add(addReq)
}

func (c *Config) LdapSendModifyRequest(addReq *ldap.ModifyRequest, domain string) error {
	conn, err := c.GetConn(domain)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Modify(addReq)
}

func (c *Config) LdapSendDeleteRequest(delReq *ldap.DelRequest, domain string) error {
	conn, err := c.GetConn(domain)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Del(delReq)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *Config) LdapCreateOrganization(group, baseDN, domain string) error {
	addReq := ldap.NewAddRequest(fmt.Sprintf("OU=%s,%s", group, baseDN), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalUnit"})
	addReq.Attribute("name", []string{group})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	return c.LdapSendAddRequest(addReq, domain)
}

func (c *Config) LdapCreateGroup(group, baseDN, domain string) error {
	addReq := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", group, baseDN), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "group"})
	addReq.Attribute("sAMAccountName", []string{group})
	addReq.Attribute("name", []string{group})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	return c.LdapSendAddRequest(addReq, domain)
}

func (c *Config) LdapCreateNewUser(name, phone, email, group, domain string) error {
	addReq := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", name, group), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "user", "person"})
	addReq.Attribute("name", []string{name})
	addReq.Attribute("sAMAccountName", []string{name})
	addReq.Attribute("userAccountControl", []string{fmt.Sprintf("%d", 0x0202)})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	addReq.Attribute("userPrincipalName", []string{email})
	addReq.Attribute("telephoneNumber", []string{phone})
	addReq.Attribute("accountExpires", []string{fmt.Sprintf("%d", 0x00000000)})
	return c.LdapSendAddRequest(addReq, domain)
}

func  (c *Config) LdapEditUser(name, phone, email, group, domain string) error{
	modReq := ldap.NewModifyRequest(fmt.Sprintf("CN=%s,%s", name, group), []ldap.Control{})
	modReq.Replace("name", []string{name})
	modReq.Replace("sAMAccountName", []string{name})
	modReq.Replace("userPrincipalName", []string{email})
	modReq.Replace("telephoneNumber", []string{phone})
	return c.LdapSendModifyRequest(modReq, domain)
}


func (c *Config) LdapChangeUserPassword(domain, user, newpassword string) error {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, err := utf16.NewEncoder().String("\"" + newpassword + "\"")
	if err != nil {
		return err
	}
	modReq := ldap.NewModifyRequest(user, []ldap.Control{})
	modReq.Replace("unicodePwd", []string{pwdEncoded})
	return c.LdapSendModifyRequest(modReq, domain)
}

func (c *Config) LdapAssignUserToGroup(targetGroup, user, domain string) error {
	modReq := ldap.NewModifyRequest(targetGroup, []ldap.Control{})
	modReq.Add("member", []string{user})
	return c.LdapSendModifyRequest(modReq, domain)
}

func (c *Config) LdapUnAssignUserToGroup(targetGroup, user, domain string) error {
	modReq := ldap.NewModifyRequest(targetGroup, []ldap.Control{})
	modReq.Delete("member", []string{user})
	return c.LdapSendModifyRequest(modReq, domain)
}


func  (c *Config) LdapUserActivate(user, domain string) error {
	modReq := ldap.NewModifyRequest(user, []ldap.Control{})
	modReq.Replace("userAccountControl", []string{fmt.Sprintf("%d", 0x0200)})
	return c.LdapSendModifyRequest(modReq, domain)
}

func (c *Config) LdapUserDeactivate(user, domain string) error {
	modReq := ldap.NewModifyRequest(user, []ldap.Control{})
	modReq.Replace("userAccountControl", []string{fmt.Sprintf("%d", 0x0202)})
	return c.LdapSendModifyRequest(modReq, domain)
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func (c *Config) LdapDeleteObject(name, domain string) error {
	delReq := ldap.NewDelRequest(name, []ldap.Control{})
	return c.LdapSendDeleteRequest(delReq, domain)
}
