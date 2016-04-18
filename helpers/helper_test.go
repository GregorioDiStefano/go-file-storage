package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	var stringList []string

	for i := 0; i < 100; i++ {
		stringList = append(stringList, RandomString(12))
	}

	for _, str := range stringList {
		assert.Equal(t, 12, len(str))
	}
}

func BenchmarkRandomString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomString(10)
	}
}
