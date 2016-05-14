package helpers

import (
	"fmt"

	"log"
	"math/rand"
	"net"
	"net/url"
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

func GetXFF(headers map[string][]string) string {
	fmt.Println("headers: ", headers)
	all, ok := headers["X-Forwarded-For"]

	if !ok {
		return ""
	}
	fmt.Println("XFF: ", all)
	possibleIP := all[0]
	if net.ParseIP(possibleIP) != nil {
		return possibleIP
	}
	return ""
}

func GetS3SignedURL(key string, filename, ip string) string {
	privKey, err := sign.LoadPEMPrivKeyFile(Config.CloudFrontPrivateKeyLocation)

	var signedURL string

	if err != nil {
		fmt.Println(err)
	}

	signer := sign.NewURLSigner(Config.CloudFrontKeyID, privKey)

	filenameEscaped := url.QueryEscape(filename)
	s3URL := fmt.Sprintf("https://%s/%s/%s", Config.CloudFrontURL, key, filenameEscaped)

	policy := &sign.Policy{
		Statements: []sign.Statement{
			{
				Resource: s3URL,
				Condition: sign.Condition{
					IPAddress:    &sign.IPAddress{SourceIP: ""},
					DateLessThan: &sign.AWSEpochTime{time.Now().Add(1 * time.Hour)},
				},
			},
		},
	}

	if len(ip) > 0 {
		signedURL, err = signer.SignWithPolicy(s3URL, policy)
	} else {
		signedURL, err = signer.Sign(s3URL, time.Now().Add(1*time.Hour))
	}

	if err != nil {
		log.Fatalf("Failed to sign url, err: %s\n", err.Error())
	}

	return signedURL
}
