package utils 

import (
	"fmt"
	"os"

	"math/rand"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
)

var Log = logrus.New()

func init() {
	rand.Seed(time.Now().UnixNano())
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	f, err := os.OpenFile(dir+"/go-file-storage.log", os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		panic(err)
	}

	Log.Out = f
	Log.Level = logrus.InfoLevel
	logrus.SetFormatter(&logrus.TextFormatter{})
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
	s3URL := fmt.Sprintf("https://%s/%s/%s", Config.GetString("cf_url"), key, filenameEscaped)

	if len(ip) > 0 {
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
	} else {
		signedURL, err = signer.Sign(s3URL, time.Now().Add(1*time.Hour))
	}

	if err != nil {
		Log.Fatalf("Failed to sign url, err: %s\n", err.Error())
	}

	return signedURL
}
