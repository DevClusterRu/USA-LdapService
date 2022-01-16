package randomHash

import (
	"math/rand"
	"strings"
	"time"
)


func init() {
	rand.Seed(time.Now().UnixNano())
}


var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates a random string of n length
func RandomString(n int) string {
	var b []rune
	for {
		b = make([]rune, n)
		for i := range b {
			b[i] = characterRunes[rand.Intn(len(characterRunes))]
		}
		if strings.ContainsAny(string(b), "1234567890"){
			break
		}
	}
	return string(b)
}