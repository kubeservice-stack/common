package utils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_MapKeys(t *testing.T) {
	assert := assert.New(t)
	var maps map[string]interface{}

	keys := Keys(maps)

	assert.Nil(keys, "not len == 0")

	maps = make(map[string]interface{})
	maps["1"] = &assert

	keys = Keys(maps)

	assert.Equal(keys, []string{"1"}, "is not equal")
}

func TestUtil_MapValues(t *testing.T) {
	assert := assert.New(t)
	var maps map[string]interface{}

	vs := Values(maps)

	assert.Nil(vs, "not len == 0")

	maps = make(map[string]interface{})
	maps["1"] = &assert

	vs = Values(maps)

	assert.Equal(vs, []interface{}{&assert}, "is not equal")
}

func TestUtil_MapSort(t *testing.T) {
	assert := assert.New(t)
	maps := make(map[string]string)
	maps["1"] = "1"
	maps["10"] = "10"
	maps["9"] = "9"
	vs := Sort(maps)

	assert.NotNil(vs, "not len == 0")
	assert.Equal(vs, map[string]string{"1": "1", "9": "9", "10": "10"}, "is not equal")

	maps["2"] = "2"

	vs = Sort(maps)

	assert.Equal(vs, map[string]string{"1": "1", "2": "2", "9": "9", "10": "10"}, "is not equal")
}

func TestUtil_MapSortKey(t *testing.T) {
	assert := assert.New(t)
	maps := make([]string, 0)
	maps = []string{"1", "10", "9"}

	vs := SortKey(maps)

	assert.NotNil(vs, "not len == 0")
	assert.Equal(vs, []string{"1", "9", "10"}, "is not equal")

	maps = []string{"1", "10", "9", "2"}

	vs = SortKey(maps)
	assert.Equal(vs, []string{"1", "2", "9", "10"}, "is not equal")

}

func TestUtil_Strings(t *testing.T) {
	assert := assert.New(t)

	retaa := make(map[string]string)

	assert.Len(retaa, 0, "dongjiang")

}

func Test_Merge(t *testing.T) {
	assert := assert.New(t)
	ret, err := Merge(map[string]string{"aa": "bb", "bb": "cc"}, map[string]string{"dd": "ee", "bb": "cc"})
	assert.Equal(ret, map[string]string{"aa": "bb", "bb": "cc", "dd": "ee"})
	assert.Nil(err)
}

func Test_ToParam(t *testing.T) {
	assert := assert.New(t)
	u, err := url.ParseQuery("https://www.baidu.com/s?ie=utf-8&f=3&rsv_bp=1&rsv_idx=1&tn=baidu&wd=docker%20environment%E8%AE%BE%E7%BD%AE%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F&fenlei=256&rsv_pq=e37e3569000003df&rsv_t=bad0PtT22fO15WdADr1HsPCYvpKKXDNWfD77qwvYtJ0JjUrF%2BvQPrcRAEqY&rqlang=cn&rsv_enter=1&rsv_dl=ts_0&rsv_sug3=10&rsv_sug1=10&rsv_sug7=100&rsv_sug2=1&rsv_btype=i&prefixsug=docker%2520env&rsp=0&inputT=4905&inputT=4905&rsv_sug4=4904")
	assert.Nil(err)
	r := ToParam(u)
	assert.Equal(r, map[string]string(map[string]string{"f": "3", "fenlei": "256", "https://www.baidu.com/s?ie": "utf-8", "inputT": "4905", "prefixsug": "docker%20env", "rqlang": "cn", "rsp": "0", "rsv_bp": "1", "rsv_btype": "i", "rsv_dl": "ts_0", "rsv_enter": "1", "rsv_idx": "1", "rsv_pq": "e37e3569000003df", "rsv_sug1": "10", "rsv_sug2": "1", "rsv_sug3": "10", "rsv_sug4": "4904", "rsv_sug7": "100", "rsv_t": "bad0PtT22fO15WdADr1HsPCYvpKKXDNWfD77qwvYtJ0JjUrF+vQPrcRAEqY", "tn": "baidu", "wd": "docker environment设置环境变量"}))
}

func Test_Strings(t *testing.T) {
	assert := assert.New(t)
	r := Strings([]interface{}{"aa", "bb", "cc", "dd", "ee"})
	assert.Equal(r, map[string]bool(map[string]bool{"aa": true, "bb": true, "cc": true, "dd": true, "ee": true}))
}

func Test_ToValues(t *testing.T) {
	assert := assert.New(t)
	r := ToValues(map[string]string{"aa": "bb"})
	assert.Equal(r, url.Values(url.Values{"aa": []string{"bb"}}))
}

func Test_ToMapStrings(t *testing.T) {
	assert := assert.New(t)
	r := ToMapStrings(map[string]interface{}{"aa": "11", "bb": "cc"})
	assert.Equal(r, map[string]string(map[string]string{"aa": "11", "bb": "cc"}))

	ret := SStrings([]string{"aaaaa", "bbbb", "c"})
	assert.Equal(ret, map[string]bool{"aaaaa": true, "bbbb": true, "c": true})
}

func Test_SliceRandList(t *testing.T) {
	assert := assert.New(t)
	aa := SliceRandList(10, 7)
	assert.Greater(len(aa), 1)
}

func Test_SliceShuffle(t *testing.T) {
	assert := assert.New(t)
	aa := SliceShuffle([]interface{}{"aa", 11})
	assert.Equal(len(aa), 2)
}

func Test_SliceDiff(t *testing.T) {
	assert := assert.New(t)
	aa := SliceDiff([]interface{}{"aa", 11}, []interface{}{"bb", 11})
	assert.Equal([]interface{}{"aa", "bb"}, aa)
}
