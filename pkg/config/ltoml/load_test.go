/*
Copyright 2022 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
