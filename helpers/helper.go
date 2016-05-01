package helpers

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length uint8) string {
	possibleCharacters := "123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	var tmp []byte

	for i := uint8(0); i < length; i++ {
		idx := rand.Intn(len(possibleCharacters))
		tmp = append(tmp, possibleCharacters[idx])
	}
	return string(tmp)
}

func IsWebBrowser(userAgent string) bool {
	cliClients := []string{"wget", "curl"}

	//check if the user agent contains a substr from cliClients
	for _, cliUA := range cliClients {
		cliUA := strings.ToLower(cliUA)
		userAgent := strings.ToLower(userAgent)
		if strings.Contains(userAgent, cliUA) {
			return false
		}
	}
	return true
}
