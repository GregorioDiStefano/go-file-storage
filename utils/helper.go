package utils

import (
	"fmt"
	"log"

	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	possibleCharacters := "123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	var tmp []byte

	for i := int(0); i < length; i++ {
		idx := rand.Intn(len(possibleCharacters))
		tmp = append(tmp, possibleCharacters[idx])
	}
	return string(tmp)
}

func IsWebBrowser(userAgent string) bool {
	cliClients := []string{"wget", "curl"}
	userAgent = strings.ToLower(userAgent)

	if userAgent == "" {
		return false
	}

	//check if the user agent contains a substr from cliClients
	for _, cliUA := range cliClients {
		cliUA := strings.ToLower(cliUA)
		if strings.Contains(userAgent, cliUA) {
			return false
		}
	}
	return true
}

func GetS3SignedURL(key string, filename, ip string) string {
	privKey, err := sign.LoadPEMPrivKeyFile(Config.GetString("aws.cf_key_location"))

	var signedURL string

	if err != nil {
		fmt.Println(err)
	}

	signer := sign.NewURLSigner(Config.GetString("aws.cf_key_id"), privKey)
	filenameEscaped := url.QueryEscape(filename)

	s3URL := fmt.Sprintf("https://%s/%s/%s", Config.GetString("aws.cf_url"), key, filenameEscaped)

	if len(ip) == 0 || strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "172.16.") || strings.HasPrefix(ip, "172.32.") {
		signedURL, err = signer.Sign(s3URL, time.Now().Add(1*time.Hour))
	} else {
		policy := &sign.Policy{
			Statements: []sign.Statement{
				{
					Resource: s3URL,
					Condition: sign.Condition{
						IPAddress:    &sign.IPAddress{SourceIP: ip},
						DateLessThan: &sign.AWSEpochTime{time.Now().Add(1 * time.Hour)},
					},
				},
			},
		}
		signedURL, err = signer.SignWithPolicy(s3URL, policy)
	}

	if err != nil {
		log.Fatalf("Failed to sign url, err: %s\n", err.Error())
	}

	return signedURL
}
