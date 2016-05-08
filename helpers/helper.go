package helpers

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
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

	if userAgent == "" {
		return false
	}

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

func GetS3SignedURL(key string, filename string) string {
	privKey, err := sign.LoadPEMPrivKeyFile("/home/greg/Desktop/keys/pk-APKAIWSICRWRZP4KRDWA.pem")

	if err != nil {
		fmt.Println(err)
	}

	signer := sign.NewURLSigner("APKAIWSICRWRZP4KRDWA", privKey)
	s3URL := fmt.Sprintf("https://dssdldrogwhnn.cloudfront.net/%s/%s", key, filename)
	signedURL, err := signer.Sign(s3URL, time.Now().Add(1*time.Hour))
	if err != nil {
		log.Fatalf("Failed to sign url, err: %s\n", err.Error())
	}

	return signedURL
}
