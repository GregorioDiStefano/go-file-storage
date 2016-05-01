package helpers

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
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
	assert.True(t, isWebBrowser("Chrome"))
	assert.True(t, isWebBrowser("chrome"))
	assert.True(t, isWebBrowser("safari"))
	assert.False(t, isWebBrowser("curl"))
	assert.False(t, isWebBrowser("wget"))
	assert.False(t, isWebBrowser("Wget"))
}

func BenchmarkRandomString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomString(10)
	}
}
