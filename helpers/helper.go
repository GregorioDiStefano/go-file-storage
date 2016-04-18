package helpers

import (
	"math/rand"
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
