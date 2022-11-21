package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetBetweenStr(t *testing.T) {
	assert := assert.New(t)

	substr := GetBetweenStr("aaaadfddweeacrttsdf", "a", "r")
	assert.Equal(substr, "aaaadfddweeac", "is true")

	substr = GetBetweenStr("aaaadfddweeacrttsdf", "d", "d")
	assert.Equal(substr, "", "is true")

	substr = GetBetweenStr("aaaadfddweeacrttsdf", "f", "t")
	assert.Equal(substr, "fddweeacr", "is true")
}

func Test_SubStr(t *testing.T) {
	assert := assert.New(t)

	substr := Substr("aaaadfddweeacrttsdf", -1, 0)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", 0, 1)
	assert.Equal(substr, "a", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -1, 1)
	assert.Equal(substr, "f", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -1, 2)
	assert.Equal(substr, "f", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2, 3)
	assert.Equal(substr, "df", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2000, 3)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2, -3)
	assert.Equal(substr, "tts", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2, -100)
	assert.Equal(substr, "aaaalfbkweeacrtts", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -1, -100)
	assert.Equal(substr, "aaaalfbkweeacrttsd", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 1000, 11111)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 4, 0)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 4, 10000)
	assert.Equal(substr, "dfddweeacrttsdf", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 4, 6)
	assert.Equal(substr, "dfddwe", "is true")
}

//func Test_interface(t *testing.T) {
//	assert := assert.New(t)

//	var tmp interface{}
//	tmp = append(tmp, "aaa")

//}
