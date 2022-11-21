package utils

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_Md5Encode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("74b87337454200d4d33f80c4663dc5e5", Md5Encode("aaaa"), "is not equal")
}

func Test_Base64Encode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("YWFhYQ==", Base64Encode("aaaa"), "is not equal")
}

func Test_urlencode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("user%3Ddongjiang%26signature%3DeBA5HZ6lccsp1jsh%252BZ7jtDFXrR61uRHHs7RV88zc2tY%253D%26expires%3D1479390425", Urlencode("user=dongjiang&signature=eBA5HZ6lccsp1jsh%2BZ7jtDFXrR61uRHHs7RV88zc2tY%3D&expires=1479390425"))
}

func Test_urldecode(t *testing.T) {
	assert := assert.New(t)
	aa, err := Urldecode("user%3Ddongjiang%26signature%3DeBA5HZ6lccsp1jsh%252BZ7jtDFXrR61uRHHs7RV88zc2tY%253D%26expires%3D1479390425")
	assert.Equal("user=dongjiang&signature=eBA5HZ6lccsp1jsh%2BZ7jtDFXrR61uRHHs7RV88zc2tY%3D&expires=1479390425", aa)
	assert.Nil(err)
}
