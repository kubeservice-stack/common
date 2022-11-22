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
	"path"
	"path/filepath"
	"testing"

	"github.com/kubeservice-stack/common/pkg/utils"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
}

var testPath = "./file"

func Test_Encode(t *testing.T) {
	_ = utils.MkDirIfNotExist(testPath)
	defer func() {
		_ = utils.RemoveDir(testPath)
	}()
	user := User{Name: "media"}
	file := path.Join(testPath, "toml")
	err := EncodeToml(file, &user)
	if err != nil {
		t.Fatal(err)
	}
	user2 := User{}
	err = DecodeToml(file, &user2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, user2)

	files, _ := utils.ListDir(testPath)
	assert.Equal(t, "toml", files[0])

	assert.NotNil(t, EncodeToml(filepath.Join(os.TempDir(), "/tmp/test.toml"), []byte{}))
}

func Test_WriteConfig(t *testing.T) {
	_ = utils.MkDirIfNotExist(testPath)
	defer func() {
		_ = utils.RemoveDir(testPath)
	}()
	assert.Nil(t, WriteConfig(path.Join(testPath, "toml"), ""))
}
