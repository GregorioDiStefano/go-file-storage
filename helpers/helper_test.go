package helpers

import (
	"sort"
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

	sort.Strings(stringList)

	for i := 0; i < len(stringList)-1; i++ {
		assert.NotEqual(t, stringList[i], stringList[i+1])
	}
}

func TestIsWebBrowser(t *testing.T) {
	assert.True(t, IsWebBrowser("Chrome"))
	assert.True(t, IsWebBrowser("chrome"))
	assert.True(t, IsWebBrowser("safari"))
	assert.False(t, IsWebBrowser("curl"))
	assert.False(t, IsWebBrowser("wget"))
	assert.False(t, IsWebBrowser("Wget"))
}

func BenchmarkRandomString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomString(10)
	}
}
