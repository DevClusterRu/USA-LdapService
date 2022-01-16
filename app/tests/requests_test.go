package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	for {

		client := &http.Client{}
		str := `{
		"user":"A.Kikos3",
			"baseDN": "CN=Users,DC=test,DC=lab",
			"command": "DropPassword"
	}`
		req, err := http.NewRequest("POST", "http://127.0.0.1:8085/ldap", bytes.NewReader([]byte(str)))
		if err != nil {
			log.Fatal("Request can`t created", err)
		}
		result, err := client.Do(req)
		if err != nil {
			log.Fatal("Request can`t sent", err)
		}
		defer req.Body.Close()
		dataTxt, err := ioutil.ReadAll(result.Body)
		if err != nil {
			log.Fatal("Request can`t read")
		}
		fmt.Println(string(dataTxt))
		time.Sleep(1*time.Second)
	}
}
