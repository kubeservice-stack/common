package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version_compare(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Version_compare("1.0", "1.0.1"), -1, "is not equals")
	assert.Equal(Version_compare("1", "1.0.1"), -1, "is not equals")
	assert.Equal(Version_compare("1.0.1", "1.0.2"), -1, "is not equals")
	assert.Equal(Version_compare("9", "139"), -1, "is not equals")
	assert.Equal(Version_compare("9", "10.111.39"), -1, "is not equals")
	assert.Equal(Version_compare("1.1.1.1.1", "1.1.1.2.1"), -1, "is not equals")
	assert.Equal(Version_compare("1.1.1.1.2", "1.1.1.2.1"), -1, "is not equals")
	assert.Equal(Version_compare("1.1.1.2.1-debug", "1.1.1.2.1"), -1, "is not equals")
	//TODO
	assert.Equal(Version_compare("1.1.1.2.2-debug", "1.1.1.2.1"), 1, "is not equals")
	assert.Equal(Version_compare("1.1.1.2.2", "1.1.1.2.1"), 1, "is not equals")
	assert.Equal(Version_compare("1.1.1.2", "1.1.1.2.1"), -1, "is not equals")
	assert.Equal(Version_compare("1.1.1.1.2-dev", "1.1.1.2.1"), -1, "is not equals")
	//TODO
	assert.Equal(Version_compare("Debug-1.1.1.1", "1.1.1.1"), -1, "is not equals")
	assert.Equal(Version_compare("debug-1.1.1.1", "1.1.1.1"), -1, "is not equals")
	assert.Equal(Version_compare("Debug-1.1.1.2", "1.1.1.1"), -1, "is not equals")
	assert.Equal(Version_compare("debug-1.1.1.2", "1.1.1.1"), -1, "is not equals")

	assert.Equal(Version_compare("44", "444"), -1, "is not equals")
	assert.Equal(Version_compare("a1", "a1.2"), -1, "is not equals")
	assert.Equal(Version_compare("a1", "a12"), -1, "is not equals")
	assert.Equal(Version_compare("aa", "444"), -1, "is not equals")
	assert.Equal(Version_compare("1.2", "1.2.3"), -1, "is not equals")
	assert.Equal(Version_compare("445.1", "444.12.1235.6667"), 1, "is not equals")
	assert.Equal(Version_compare("4.4.1", "125"), -1, "is not equals")

	//TODO maybe bug
	assert.Equal(Version_compare("aa", "aaaaa"), 0, "is not equals")
}
