package ltoml

import (
	"os"
	"testing"

	"github.com/kubeservice-stack/common/pkg/utils"

	"github.com/stretchr/testify/assert"
)

type TestCfg struct {
	Path string `toml:"path"`
}

var cfgFile = "./test.test"
var defaultCfgFile = "./test.test"

func Test_LoadConfig(t *testing.T) {
	defer func() {
		_ = utils.RemoveDir(cfgFile)
	}()
	assert.NotNil(t, LoadConfig(cfgFile, defaultCfgFile, &TestCfg{}))

	f, err := os.Create(cfgFile)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("dongjiang test")
	assert.NotNil(t, LoadConfig(cfgFile, defaultCfgFile, &TestCfg{}))

	_ = EncodeToml(cfgFile, &TestCfg{Path: "/data/path"})
	cfg := TestCfg{}
	err = LoadConfig(cfgFile, defaultCfgFile, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)

	err = LoadConfig("", defaultCfgFile, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)
}
